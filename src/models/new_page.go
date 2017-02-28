package models

import (
	"strings"
	"time"
	"models/status"
	"models/visibility"
)

type NewPage struct {
	Title      string
	Content    string
	Template   string
	Status     status.Status
	ScheduleAt time.Time `bson:"scheduleat,omitempty"`
	Author     *UserPrincipal
}

type NewPageRequest struct {
	Title      string
	Content    string
	Template   string
	Status     status.Status
	Visibility visibility.Visibility
	ScheduleAt time.Time
}

func (p NewPageRequest) Validate() (bool, *Errs) {
	var errors Errs
	title := strings.TrimSpace(p.Title)
	content := strings.TrimSpace(p.Content)

	switch {
	case IsEmpty(title):
		errors.AddEmptyError("Title")

	case LengthNotBetween(title, 2, 255):
		errors.AddMaxLengthError(255, "Title", p.Title)

	case ContainsNotAllowed(p.Title, PatternAlphaNumericPun):
		errors.AddNotAllowedError("Title", p.Title, PatternAlphaNumericPun)

	case IsEmpty(p.Template):
		errors.AddEmptyError("Template")

	case LengthExceedsMax(content, 2048):
		errors.AddMaxLengthError(2048, "Content", p.Content)

	case IsEmpty(title):
		errors.AddEmptyError("Title")

	case status.IsNotValid(p.Status):
		errors.AddInvalidError("Status", p.Status.String())

	case visibility.IsNotValid(p.Visibility):
		errors.AddInvalidError("Visibility", p.Visibility.String())
	}

	if len(errors) > 0 {
		return false, &errors
	}

	return true, nil
}
