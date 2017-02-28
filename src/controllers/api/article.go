package api

import (
	"time"
)

type Article struct {
	Title      string
	Slug       string
	Content    string
	Excerpt    string
	Author     struct {
		           FullName          string
		           NickName          string
		           Website           string
		           Bio               string
		           ProfilePictureUrl string
	           }
	Categories []struct {
		Name string
		Slug string
	}
	Taxonomies []struct {
		Name string
		Slug string
	}
	UpdatedAt  time.Time
}
