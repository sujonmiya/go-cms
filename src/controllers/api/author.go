package api

import (
	"net/url"
	"service/app/models/gender"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type AuthorResponse struct {
	FullName       string
	NickName       string
	Gender         gender.Gender
	Email          string
	Website        url.URL
	Bio            string
	ProfilePicture bson.ObjectId
	Articles       []Article
	UpdatedAt      time.Time
}
