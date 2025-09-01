package budget

type Status string

const (
	OnTrack  Status = "ontrack"  // <= 50% usage
	Warning  Status = "warning"  // >= 50% usage
	Caution  Status = "caution"  // >= 10% usage
	Critical Status = "critical" // >= 90% usage
)
