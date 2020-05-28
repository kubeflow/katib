package trial

const (
	DefaultJobKind = "Job"

	// For trials
	TrialCreatedReason            = "TrialCreated"
	TrialRunningReason            = "TrialRunning"
	TrialSucceededReason          = "TrialSucceeded"
	TrialMetricsUnavailableReason = "MetricsUnavailable"
	TrialFailedReason             = "TrialFailed"
	TrialKilledReason             = "TrialKilled"

	// For Jobs
	JobCreatedReason            = "JobCreated"
	JobDeletedReason            = "JobDeleted"
	JobSucceededReason          = "JobSucceeded"
	JobMetricsUnavailableReason = "MetricsUnavailable"
	JobFailedReason             = "JobFailed"
	JobRunningReason            = "JobRunning"
	ReconcileFailedReason       = "ReconcileFailed"
)
