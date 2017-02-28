package models

import (
	"strings"
	"time"
	"models/status"
	"models/visibility"
	"github.com/markbates/pop/nulls"
)

type Article struct {
	Model
	Title        string `sql:"not null"`
	Slug         string `sql:"not null;unique;index"`
	Content      string `sql:"type:TEXT"`
	Excerpt      string
	Picture      *Picture
	PictureID    nulls.UInt32
	Status       string `sql:"not null;index;size:25"`
	AuthorID     nulls.UInt32 `sql:"not null;index"`
	LastEditorID nulls.UInt32
}

type UpdateArticleRequest struct {
	Title         string
	Content       string
	Excerpt       string
	FeaturedImage uint32
	Categories    []uint32
	Tags          []uint32
	Status        status.Status
	ScheduleAt    time.Time `json:",omitempty"`
}

func (a UpdateArticleRequest) Validate() (bool, *Errs) {
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
	}

	if len(errors) > 0 {
		return false, &errors
	}

	return true, nil
}

type UpdateArticle struct {
	Title      string                `bson:"title"`
	Slug       string                `bson:"slug"`
	Content    string                `bson:"content"`
	Metadata   Metadata              `bson:"metadata"`
	Categories []Taxonomy            `bson:"categories"`
	Tags       []Taxonomy            `bson:"tags"`
	Status     status.Status         `bson:"status"`
	Visibility visibility.Visibility `bson:"visibility"`
	ScheduleAt time.Time             `bson:"scheduleat"`
	Editor     *UserPrincipal        `bson:"editor"`
	UpdatedAt  time.Time
}

/*func (a Article) Excerpt() string {
	s := strings.Fields(a.Content)
	if len(s) > 15 {
		return strings.Join(s[:15], " ")
	}

	return a.Content
}

func (a *Article) AddCategory(category *Category) {
	c := Category{
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
	}

	a.Categories = append(a.Categories, c)
}

func (a *Article) AddCategories(categories []*Category) {
	for _, category := range categories {
		a.AddCategory(category)
	}
}

func (a *Article) AddTag(tag *Tag) {
	t := Taxonomy{
		Name:        tag.Name,
		Slug:        tag.Slug,
		Description: tag.Description,
	}

	a.Tags = append(a.Tags, t)
}

func (a *Article) AddTags(tags []*Tag) {
	for _, tag := range tags {
		a.AddTag(tag)
	}
}

func (a *Article) SetAuthor(user *UserPrincipal) {
	a.Author = user.Id
}

func (a *Article) SetEditor(user *UserPrincipal) {
	a.LastEditor = user.Id
}*/
