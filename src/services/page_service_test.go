package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"models"
	"utils/seed"
	"repository"
)

var (
	pageService *PageService
	page *models.Page
)

func init() {
	_service := NewPageService()
	pageService = _service
}

func TestPageService(t *testing.T) {
	a := assert.New(t)
	a.NotNil(pageService)
}

func TestPageService_SavePage(t *testing.T) {
	_page, err := pageService.SavePage(seed.NewPage())
	a := assert.New(t)
	a.NoError(err)
	a.NotNil(_page)
	page = _page
}

func TestPageService_GetPageBySlug(t *testing.T) {
	p, err := pageService.GetPageBySlug(page.Slug)
	a := assert.New(t)
	a.NoError(err)
	a.NotNil(p)
}

func TestPageService_GetPageByID(t *testing.T) {
	ass := assert.New(t)
	page, err := pageService.GetPageByID(1)
	ass.NoError(err)
	ass.NotNil(page)
}

func TestPageService_GetPages(t *testing.T) {
	pages, err := pageService.GetPages(repository.NewDefaultQuery())
	a := assert.New(t)
	a.NoError(err)
	a.NotEmpty(pages)
}

func TestPageService_GetPagesByQuery(t *testing.T) {
	ass := assert.New(t)
	pages, err := pageService.GetPagesByQuery(repository.NewDefaultQuery())
	ass.NoError(err)
	ass.NotEmpty(pages)
	ass.Len(pages, 10)
}

func TestPageService_GetPage(t *testing.T) {
	ass := assert.New(t)
	page , err := pageService.GetPage("aspernatur-exercitationem-est-quo-commodi")
	ass.NoError(err)
	ass.NotNil(page)
}


