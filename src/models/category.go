package models

import (
	"time"

	"utils"
	"github.com/markbates/pop/nulls"
)

type Category struct {
	Model
	Name         string `sql:"not null"`
	Slug         string `sql:"not null;unique;index"`
	Description  string `sql:"type:TEXT"`
	ParentID     nulls.UInt32       `json:",omitempty" sql:"index"`
	AuthorID     nulls.UInt32 `sql:"not null;index"`
	LastEditorID nulls.UInt32
}

type UpdateCategoryRequest struct {
	Name        string
	Description string        `json:",omitempty"`
	Parent      string `json:",omitempty"`
}

func (c UpdateCategoryRequest) Validate() (bool, *Errs) {
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

type UpdateCategory struct {
	Name        string
	Slug        string
	Description string         `bson:"description,omitempty"`
	Parent      Taxonomy       `bson:",omitempty"`
	Children    []Taxonomy     `bson:",omitempty"`
	Editor      *UserPrincipal `bson:"editor,omitempty"`
	UpdatedAt   time.Time
}

func (c *Category) SetParent(parent Taxonomy) {
	c.ParentID = utils.ToUInt32(parent.ID)
}