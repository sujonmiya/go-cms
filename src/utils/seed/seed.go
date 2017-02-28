package seed

import (
	"github.com/icrowley/fake"
	"models"
	"models/status"
	"time"
	"strconv"
	"github.com/enodata/faker"
	"strings"
	"io/ioutil"
	"net/http"
	"models/roles"
	"models/capabilities"
	"log"
	"fmt"
)

func NewAdministrator() *models.NewUser {
	return &models.NewUser{
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		Email:     fake.EmailAddress(),
		Password:  []byte(fake.SimplePassword()),
		Role:      roles.Administrator,
	}
}

func NewAuthor() models.NewUser {
	return models.NewUser{
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		Email:     fake.EmailAddress(),
		Password:  []byte(fake.SimplePassword()),
		Role:      roles.Author,
	}
}

func NewSubscriber() models.NewUser {
	return models.NewUser{
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		Email:     fake.EmailAddress(),
		Password:  []byte(fake.SimplePassword()),
		Role:      roles.Subscriber,
	}
}

func NewCategory() *models.NewCategory {
	c := &models.NewCategory{
		Name:        fake.Word(),
		Description: fake.Paragraphs(),
		Author: NewUserPrincipal(),
	}

	return c
}

func NewTaxonomy() *models.NewTaxonomy {
	return &models.NewTaxonomy{
		Name:        fake.Word(),
		Description: fake.Paragraphs(),
		Author: NewUserPrincipal(),
	}
}

func NewArticle() *models.NewArticle {
	return &models.NewArticle{
		Title:      fake.Title(),
		Content:    fake.Paragraphs(),
		Excerpt: fake.WordsN(10),
		FeaturedImage: NewPicture(),
		Status:     status.Draft,
		ScheduleAt: time.Now(),
	}
}

func NewUserPrincipal() *models.UserPrincipal{
	return &models.UserPrincipal{
		ID:1,
		Role:roles.Administrator,
		Capabilities:capabilities.Capabilities(),
	}
}

func NewPage() *models.NewPage {
	p := &models.NewPage{
		Title:      fake.Title(),
		Content:    fake.Paragraphs(),
		Template: fake.Word(),
		Status:     status.Draft,
		ScheduleAt: time.Now().Add(time.Hour * 12),
	}

	p.Author = NewUserPrincipal()
	return p
}

func Picture() *models.Picture {
	w, _ := strconv.Atoi(fake.DigitsN(3))
	h, _ := strconv.Atoi(fake.DigitsN(3))
	s, _ := strconv.Atoi(fake.DigitsN(5))
	p := &models.Picture{
		Name:      strings.Replace(fake.Title(), " ", "_", -1),
		Width:    w,
		Height:     h,
		Size: s,
		Caption: fake.Words(),
		AltText: fake.Words(),
		MimeType:"image/png",
		Url: "http://loremflickr.com/320/240",
	}
	id, _ := strconv.Atoi(faker.Number().Number(5))
	p.ID = uint32(id)
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt

	return p
}

func NewPicture() *models.NewPicture {
	filename := `C:\Users\Sujon\Documents\Projects\Contetto\cms\uploads\2016\12\2\d6dww7lmku_600x400_15e269cfb913bec01fbbf8a65c65415395e0fce83a20651b0729641211bf0c9d.jpg`
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Could not read request body: %v", err)
		return nil
	}

	p := &models.NewPicture{
		Name:      fmt.Sprintf("%s-%s.jpg", fake.Word(), fake.Word()),
		Data: data,
		Caption: fake.Words(),
		AltText: fake.Words(),
		MimeType:http.DetectContentType(data),
	}

	return p
}

func NewLogin() models.Login {
	return models.Login{
		Email:    fake.EmailAddress(),
		Password: []byte(fake.SimplePassword()),
	}
}
