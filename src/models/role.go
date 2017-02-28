package models

type Role struct {
	Model
	Name        string `sql:"not null;unique;index"`
	Description string        `json:"description,omitempty"`
}
