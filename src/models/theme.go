package models

type Theme struct {
	Model
	Name        string
	Version     float32
	Description string `yaml:",omitempty"`
	Author      struct {
		Name    string
		Website string
	}
	Templates []string `yaml:",omitempty"`
}
