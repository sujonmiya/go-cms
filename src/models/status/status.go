package status

//go:generate jsonenums -type=Gender
//go:generate stringer -type=Status
type Status int64

const (
	Draft Status = iota
	PendingReview
	Published
)

func Statuses() []string {
	var statuses []string
	for _, s := range _StatusValueToName {
		statuses = append(statuses, s)
	}

	return statuses
}

func IsNotValid(status Status) bool {
	for _, s := range Statuses() {
		if s != status.String() {
			return false
		}
	}

	return true
}
