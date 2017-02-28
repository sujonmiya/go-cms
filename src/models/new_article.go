package models

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
	"models/status"
	"models/visibility"
)

type NewArticle struct {
	Title         string                `bson:"title"`
	Content       string                `bson:"content"`
	Excerpt       string
	Categories    []uint32            `bson:"categories"`
	Taxonomies    []uint32            `bson:"tags"`
	FeaturedImage *NewPicture
	Status        status.Status         `bson:"status"`
	ScheduleAt    time.Time      `bson:"scheduleat"`
	Author        *UserPrincipal
}

func (art *NewArticle) HasFeaturedImage() bool {
	return art.FeaturedImage != nil
}

type NewArticleRequest struct {
	Title      string
	Content    string
	Categories []bson.ObjectId
	Tags       []bson.ObjectId `json:",omitempty"`
	Metadata   Metadata        `json:",omitempty"`
	Status     status.Status
	Visibility visibility.Visibility
	ScheduleAt time.Time `json:",omitempty"`
}

func (a NewArticleRequest) Validate() (bool, *Errs) {
	var errors Errs
	title := strings.TrimSpace(a.Title)
	content := strings.TrimSpace(a.Content)

	switch {
	case IsEmpty(title):
		errors.AddEmptyError("Title")

	case LengthNotBetween(title, 2, 255):
		errors.AddMaxLengthError(255, "Title", a.Title)

	case ContainsNotAllowed(a.Title, PatternAlphaNumericPun):
		errors.AddNotAllowedError("Title", a.Title, PatternAlphaNumericPun)

	case IsEmpty(content):
		errors.AddEmptyError("Content")

	case LengthExceedsMax(content, 2048):
		errors.AddMaxLengthError(2048, "Content", a.Content)

	case len(a.Categories) == 0:
		errors.AddEmptyError("Categories")

	case status.IsNotValid(a.Status):
		errors.AddInvalidError("Status", a.Status.String())

	case visibility.IsNotValid(a.Visibility):
		errors.AddInvalidError("Visibility", a.Visibility.String())
	}

	if len(errors) > 0 {
		return false, &errors
	}

	return true, nil
}
