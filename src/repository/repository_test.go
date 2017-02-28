package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"models"
	"github.com/icrowley/fake"
	"strings"
	"utils/seed"
)

var (
	repo *Repository
)

func init() {
	_repo := NewRepo()
	repo = _repo
}

func TestNewRepository(t *testing.T) {
	a := assert.New(t)
	a.NotNil(repo)
}

func TestRepository_Save(t *testing.T) {
	a := assert.New(t)
	err := repo.Save(seed.Picture())
	a.NoError(err)
}

func TestRepository_FormatSort(t *testing.T) {
	ass := assert.New(t)
	sort := formatSort("   created_at   ,   slug    ")
	ass.Equal("created_at slug", sort)
}

func TestRepository_Find(t *testing.T) {
	a := assert.New(t)
	var pics []*models.Picture
	err := repo.Find(&pics)
	a.NoError(err)
	a.NotEmpty(pics)
}

func TestRepository_FindByID(t *testing.T) {
	ass := assert.New(t)
	var pic models.Picture
	err := repo.FindByID(1, &pic)
	ass.NoError(err)
	ass.NotNil(pic)
}

func TestRepository_FindByQuery(t *testing.T) {
	a := assert.New(t)
	var pics []*models.Picture
	err := repo.FindByQuery(seed.NewDefaultQuery(), &pics)
	a.NoError(err)
	a.NotEmpty(pics)
}

func TestRepository_FindByQueryAndFilter(t *testing.T) {
	a := assert.New(t)
	mime := "image/png"
	filter := &models.Picture{
		MimeType: mime,
		Size:462,
	}
	var pics []*models.Picture
	err := repo.FindByQueryAndFilter(seed.NewDefaultQuery(), filter, &pics)
	a.NoError(err)

	for _, p := range pics {
		a.Equal(mime, p.MimeType)
	}
}

func TestRepository_FindOne(t *testing.T) {
	a := assert.New(t)
	mime := "image/png"
	filter := &models.Picture{
		MimeType: mime,
		Width:343,
	}
	var pic models.Picture
	err := repo.FindOne(filter, &pic)
	a.NoError(err)
	a.NotNil(pic)
}

func TestRepository_Update(t *testing.T) {
	a := assert.New(t)
	filter := &models.Picture{Model:models.Model{ID:2}}
	picture := &models.Picture{
		Name:strings.Replace(fake.ProductName(), " ", "_", -1),
	}

	err := repo.Update(filter, picture)
	a.NoError(err)
}

func TestRepository_Delete(t *testing.T) {
	a := assert.New(t)
	filter := &models.Picture{Model:models.Model{ID:1}}
	err := repo.Delete(filter)
	a.NoError(err)
}