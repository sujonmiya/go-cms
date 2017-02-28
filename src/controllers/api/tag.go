package api

import (
	"net/mail"
	"net/url"
	"time"
)

type TagResponse struct {
	Name        string
	Slug        string
	Description string
	Articles    []struct {
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
		Categories []struct {
			Name        string
			Slug        string
			Description string
		}
	}
	UpdatedAt time.Time
}
