package txindex

import "github.com/go-kit/kit/metrics"

const MetricsSubsystem = "tx_indexer"

//go:generate go run ../../scripts/metricsgen -struct=Metrics

// Metrics contains timing and throughput metrics for the asynchronous indexer.
type Metrics struct {
	// Time spent receiving all transaction events for one block.
	GatherEventsSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent writing one block to the block index.
	BlockIndexSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent writing one transaction batch to the transaction index.
	TxIndexSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// End-to-end time spent processing one indexed block.
	BlockTotalSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Number of transactions handed to the transaction index.
	IndexedTxsTotal metrics.Counter
}
