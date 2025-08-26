package budget

type Status string

const (
	Healthty   Status = "healthty"
	Completed  Status = "completed"
	OverBudget Status = "overbudget"
	Suspended  Status = "suspended"
)
