package api

import (
	"net/mail"
	"net/url"
	"time"
)

type CategoryResponse struct {
	Name        string
	Slug        string
	Description string
	Parent      struct {
		Name        string
		Slug        string
		Description string
	}
	Children []struct {
		Name        string
		Slug        string
		Description string
	}
	Articles []struct {
		Title   string
		Slug    string
		Content string
		Author  struct {
			FullName       string
			NickName       string
			Email          mail.Address
			Website        url.URL
			Bio            string
			ProfilePicture url.URL
		}
		Tags []struct {
			Name string
			Slug string
		}
	}
	UpdatedAt time.Time
}

type Categories []CategoryResponse
