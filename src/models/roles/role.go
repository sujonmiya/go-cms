package roles

//go:generate jsonenums -type=Role
//go:generate stringer -type=Role
type Role uint8

const (
	_ = iota
	Administrator Role = iota
	Author
	Editor
	Subscriber
)

func Roles() []Role {
	var roles []Role
	for _, r := range _RoleNameToValue {
		roles = append(roles, r)
	}

	return roles
}

func IsNotValid(role Role) bool {
	for _, r := range Roles() {
		if r.String() != role.String() {
			return false
		}
	}

	return true
}
