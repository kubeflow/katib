package trial

const (
	DefaultJobKind = "Job"

	// For trials
	TrialCreatedReason   = "TrialCreated"
	TrialRunningReason   = "TrialRunning"
	TrialSucceededReason = "TrialSucceeded"
	TrialFailedReason    = "TrialFailed"
	TrialKilledReason    = "TrialKilled"

	// For Jobs
	JobCreatedReason         = "JobCreated"
	JobDeletedReason         = "JobDeleted"
	JobSucceededReason       = "JobSucceeded"
	MetricsUnavailableReason = "MetricsUnavailable"
	JobFailedReason          = "JobFailed"
)
