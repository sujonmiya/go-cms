package models

type NewCategory struct {
	Name        string
	Description string
	Parent      uint32
	Author      *UserPrincipal
}

type NewCategoryRequest struct {
	Name        string
	Description string        `json:",omitempty"`
	Parent      string `json:",omitempty"`
}

func (c NewCategoryRequest) Validate() (bool, *Errs) {
	var errors Errs

	switch {
	case IsEmpty(c.Name):
		errors.AddEmptyError("Name")

	case LengthNotBetween(c.Name, 2, 25):
		errors.AddMaxLengthError(25, "Name", c.Name)

	case ContainsNotAllowed(c.Name, PatternAlphaSpace):
		errors.AddNotAllowedError("Name", c.Name, PatternAlphaSpace)
	}

	if len(errors) > 0 {
		return false, &errors
	}

	return true, nil
}
