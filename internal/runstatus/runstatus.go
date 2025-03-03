// Package runstatus provides run statuses.
//
// NOTE: placed in separate package from `run` to avoid import cycles.
package runstatus

// Status represents a run state.
type Status string

const (
	// List all available run statuses supported in OTF.
	Applied            Status = "applied"
	ApplyQueued        Status = "apply_queued"
	Applying           Status = "applying"
	Canceled           Status = "canceled"
	Confirmed          Status = "confirmed"
	Discarded          Status = "discarded"
	Errored            Status = "errored"
	ForceCanceled      Status = "force_canceled"
	Pending            Status = "pending"
	PlanQueued         Status = "plan_queued"
	Planned            Status = "planned"
	PlannedAndFinished Status = "planned_and_finished"
	Planning           Status = "planning"

	// OTF doesn't support cost estimation but go-tfe API tests expect this
	// status so it is included expressly to pass the tests.
	CostEstimated Status = "cost_estimated"
)

func (s Status) String() string { return string(s) }
