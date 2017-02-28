package models

import (
	"github.com/markbates/pop/nulls"
	"github.com/go-ozzo/ozzo-validation"
)

type Picture struct {
	Model
	Name         string `sql:"not null"`
	Width        int `sql:"not null"`
	Height       int `sql:"not null"`
	Size         int `sql:"not null"`
	Caption      string
	AltText      string
	MimeType     string `sql:"not null"`
	Url          string `sql:"not null"`
	AuthorID     nulls.UInt32 `sql:"not null;index"`
	LastEditorID nulls.UInt32
}

type NewPicture struct {
	Name     string
	Data     []byte
	Caption  string
	AltText  string
	MimeType string
}

func (p NewPicture) Validate() error {
	rules := validation.StructRules{}
	rules.Add("Name", validation.Required, validation.Length(6, 255))

	return rules.Validate(p)
}