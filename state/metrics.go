package state

import (
	"github.com/go-kit/kit/metrics"
)

const (
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this
	// package.
	MetricsSubsystem = "state"
)

//go:generate go run ../scripts/metricsgen -struct=Metrics

// Metrics contains metrics exposed by this package.
type Metrics struct {
	// Time spent in the complete ApplyVerifiedBlock state execution path.
	ApplyBlockSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent executing FinalizeBlock through the consensus ABCI connection.
	FinalizeBlockSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent persisting the FinalizeBlock response.
	SaveFinalizeBlockResponseSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent deriving the next in-memory consensus state.
	UpdateStateSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent in BlockExecutor.Commit, including mempool synchronization.
	BlockCommitSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent waiting to acquire the mempool lock during block commit.
	MempoolLockWaitSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time the mempool lock is held during block commit.
	MempoolLockHeldSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent flushing outstanding mempool ABCI requests before app commit.
	FlushAppConnSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent executing the ABCI Commit call.
	AppCommitSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent updating the mempool after app commit.
	MempoolUpdateSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent updating the evidence pool after app commit.
	EvidenceUpdateSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent saving the post-commit consensus state.
	StateSaveSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent publishing block and transaction events.
	FireEventsSeconds metrics.Histogram `metrics_buckettype:"exprange" metrics_bucketsizes:"0.0001, 10, 16"`

	// Time spent processing FinalizeBlock
	BlockProcessingTime metrics.Histogram `metrics_buckettype:"lin" metrics_bucketsizes:"1, 10, 10"`

	// ConsensusParamUpdates is the total number of times the application has
	// updated the consensus params since process start.
	//metrics:Number of consensus parameter updates returned by the application since process start.
	ConsensusParamUpdates metrics.Counter

	// ValidatorSetUpdates is the total number of times the application has
	// updated the validator set since process start.
	//metrics:Number of validator set updates returned by the application since process start.
	ValidatorSetUpdates metrics.Counter

	// The number of transactions rejected by the application.
	RejectedTransactions metrics.Counter

	// The number of transactions processed by the application.
	ProcessedTransactions metrics.Counter
}
