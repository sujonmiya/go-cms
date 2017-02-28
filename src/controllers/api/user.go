package api

import (
	"net/url"
	"time"
)

type UserResponse struct {
	FullName string
	NickName string
	Email    string
	Website  url.URL
	Bio      string
	Articles []struct {
		Title      string
		Slug       string
		Content    string
		Categories []struct {
			Name        string
			Slug        string
			Description string
		}
		Tags []struct {
			Name        string
			Slug        string
			Description string
		}
		UpdatedAt time.Time
	}
}
