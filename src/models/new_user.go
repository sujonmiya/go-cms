package models

import (
	"models/roles"
	"models/capabilities"
)

type NewUser struct {
	FirstName string
	LastName  string
	Email     string
	Password  []byte
	Role      roles.Role
	CreatedBy *UserPrincipal
}

type NewUserRequest struct {
	FirstName string
	LastName  string `json:",omitempty"`
	Email     string
	Password  []byte
	Role      roles.Role
}

type UserPrincipal struct {
	ID           uint32
	Role         roles.Role
	Capabilities []capabilities.Capability
}

func (u NewUserRequest) Validate() (bool, *Errs) {
	var errors Errs
	p := string(u.Password)

	switch {
	case IsEmpty(u.FirstName):
		errors.AddEmptyError("FirstName")

	case LengthNotBetween(u.FirstName, 2, 64):
		errors.AddLengthBetweenError(2, 64, "FirstName", u.FirstName)

	case ContainsNotAllowed(u.FirstName, PatternAlphaSpace):
		errors.AddNotAllowedError("FirstName", u.FirstName, PatternAlphaSpace)

	case LengthExceedsMax(u.LastName, 64):
		errors.AddMaxLengthError(64, "LastName", u.LastName)

	case IsNotEmail(u.Email):
		errors.AddInvalidError("Email", u.Email)

	case IsEmpty(p):
		errors.AddEmptyError("Password")

	case LengthNotBetween(p, 6, 16):
		errors.AddLengthBetweenError(6, 16, "Password", u.FirstName)

	case roles.IsNotValid(u.Role):
		errors.AddInvalidError("Role", u.Role.String())
	}

	if len(errors) > 0 {
		return false, &errors
	}

	return true, nil
}
