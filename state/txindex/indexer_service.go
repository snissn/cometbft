package txindex

import (
	"context"
	"time"

	"github.com/cometbft/cometbft/libs/service"
	"github.com/cometbft/cometbft/state/indexer"
	"github.com/cometbft/cometbft/types"
)

// XXX/TODO: These types should be moved to the indexer package.

const (
	subscriber = "IndexerService"
)

// IndexerService connects event bus, transaction and block indexers together in
// order to index transactions and blocks coming from the event bus.
type IndexerService struct {
	service.BaseService

	txIdxr           TxIndexer
	blockIdxr        indexer.BlockIndexer
	eventBus         *types.EventBus
	terminateOnError bool
	metrics          *Metrics
}

type IndexerServiceOption func(*IndexerService)

func IndexerServiceWithMetrics(metrics *Metrics) IndexerServiceOption {
	return func(service *IndexerService) {
		service.metrics = metrics
	}
}

// NewIndexerService returns a new service instance.
func NewIndexerService(
	txIdxr TxIndexer,
	blockIdxr indexer.BlockIndexer,
	eventBus *types.EventBus,
	terminateOnError bool,
	options ...IndexerServiceOption,
) *IndexerService {

	is := &IndexerService{
		txIdxr:           txIdxr,
		blockIdxr:        blockIdxr,
		eventBus:         eventBus,
		terminateOnError: terminateOnError,
		metrics:          NopMetrics(),
	}
	for _, option := range options {
		option(is)
	}
	is.BaseService = *service.NewBaseService(nil, "IndexerService", is)
	return is
}

// OnStart implements service.Service by subscribing for all transactions
// and indexing them by events.
func (is *IndexerService) OnStart() error {
	// Use SubscribeUnbuffered here to ensure both subscriptions does not get
	// canceled due to not pulling messages fast enough. Cause this might
	// sometimes happen when there are no other subscribers.
	blockSub, err := is.eventBus.SubscribeUnbuffered(
		context.Background(),
		subscriber,
		types.EventQueryNewBlockEvents)
	if err != nil {
		return err
	}

	txsSub, err := is.eventBus.SubscribeUnbuffered(context.Background(), subscriber, types.EventQueryTx)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-blockSub.Canceled():
				return
			case msg := <-blockSub.Out():
				blockStart := time.Now()
				eventNewBlockEvents := msg.Data().(types.EventDataNewBlockEvents)
				height := eventNewBlockEvents.Height
				numTxs := eventNewBlockEvents.NumTxs

				batch := NewBatch(numTxs)

				gatherStart := time.Now()
				for i := int64(0); i < numTxs; i++ {
					msg2 := <-txsSub.Out()
					txResult := msg2.Data().(types.EventDataTx).TxResult

					if err = batch.Add(&txResult); err != nil {
						is.Logger.Error(
							"failed to add tx to batch",
							"height", height,
							"index", txResult.Index,
							"err", err,
						)

						if is.terminateOnError {
							if err := is.Stop(); err != nil {
								is.Logger.Error("failed to stop", "err", err)
							}
							return
						}
					}
				}
				is.metrics.GatherEventsSeconds.Observe(time.Since(gatherStart).Seconds())

				blockIndexStart := time.Now()
				err = is.blockIdxr.Index(eventNewBlockEvents)
				is.metrics.BlockIndexSeconds.Observe(time.Since(blockIndexStart).Seconds())
				if err != nil {
					is.Logger.Error("failed to index block", "height", height, "err", err)
					if is.terminateOnError {
						if err := is.Stop(); err != nil {
							is.Logger.Error("failed to stop", "err", err)
						}
						return
					}
				} else {
					is.Logger.Info("indexed block events", "height", height)
				}

				txIndexStart := time.Now()
				err = is.txIdxr.AddBatch(batch)
				is.metrics.TxIndexSeconds.Observe(time.Since(txIndexStart).Seconds())
				if err != nil {
					is.Logger.Error("failed to index block txs", "height", height, "err", err)
					if is.terminateOnError {
						if err := is.Stop(); err != nil {
							is.Logger.Error("failed to stop", "err", err)
						}
						return
					}
				} else {
					is.metrics.IndexedTxsTotal.Add(float64(numTxs))
					is.Logger.Debug("indexed transactions", "height", height, "num_txs", numTxs)
				}
				is.metrics.BlockTotalSeconds.Observe(time.Since(blockStart).Seconds())
			}
		}
	}()
	return nil
}

// OnStop implements service.Service by unsubscribing from all transactions.
func (is *IndexerService) OnStop() {
	if is.eventBus.IsRunning() {
		_ = is.eventBus.UnsubscribeAll(context.Background(), subscriber)
	}
}
