package visibility

//go:generate jsonenums -type=Gender
//go:generate stringer -type=Visibility
type Visibility uint8

const (
	Public Visibility = iota
	Private
)

func Visibilities() []string {
	var visibilities []string
	for _, v := range _VisibilityValueToName {
		visibilities = append(visibilities, v)
	}

	return visibilities
}

func IsNotValid(visibility Visibility) bool {
	for _, v := range Visibilities() {
		if v != visibility.String() {
			return false
		}
	}

	return true
}
