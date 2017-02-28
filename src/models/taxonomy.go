package models

import "github.com/markbates/pop/nulls"

type Taxonomy struct {
	Model
	Name         string `sql:"not null"`
	Slug         string `sql:"not null;unique;index"`
	Description  string `sql:"type:TEXT"`
	ParentID     nulls.UInt32 `sql:"index"`
	AuthorID     nulls.UInt32 `sql:"not null;index"`
	LastEditorID nulls.UInt32
}

type NewTaxonomy struct {
	Name        string
	Description string
	Parent      uint32
	Author      *UserPrincipal
}
