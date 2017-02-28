package models

import (
	"strings"
	"time"
	"models/status"
	"models/visibility"
	"github.com/markbates/pop/nulls"
)

type Page struct {
	Model
	Title        string `sql:"not null"`
	Slug         string `sql:"not null;unique;index"`
	Content      string `sql:"type:TEXT"`
	Template     string
	Status       string `sql:"not null;index;size:25"`
	AuthorID     nulls.UInt32 `sql:"not null;index"`
	LastEditorID nulls.UInt32
}

type UpdatePage struct {
	Title      string
	Slug       string
	Content    string
	Template   string
	Status     status.Status
	Visibility visibility.Visibility
	ScheduleAt time.Time
	Editor     *UserPrincipal
	UpdatedAt  time.Time
}

type UpdatePageRequest struct {
	Title      string
	Content    string
	Template   string
	Status     status.Status
	Visibility visibility.Visibility
	ScheduleAt time.Time
}

func (p UpdatePageRequest) Validate() (bool, *Errs) {
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

	case IsEmpty(content):
		errors.AddEmptyError("Content")

	case LengthExceedsMax(content, 2048):
		errors.AddMaxLengthError(2048, "Content", p.Content)

	case IsEmpty(p.Template):
		errors.AddEmptyError("Template")

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
