package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"models"
	"utils/seed"
	"repository"
)

var (
	categoryService *CategoryService
	category *models.Category
)

func init() {
	_service := NewCategoryService()
	categoryService = _service
}

func TestNewCategoryService(t *testing.T) {
	a := assert.New(t)
	a.NotNil(categoryService)
}

func TestCategoryService_SaveCategory(t *testing.T) {
	a := assert.New(t)
	c := seed.NewCategory()
	c.Parent = 2
	_category, err := categoryService.SaveCategory(c)
	a.NoError(err)
	a.NotNil(_category)
	category = _category
}

func TestCategoryService_GetCategory(t *testing.T) {
	ass := assert.New(t)
	cat, err := categoryService.GetCategoryByID(1)
	ass.NoError(err)
	ass.NotNil(cat)
}

func TestCategoryService_GetCategoryBySlug(t *testing.T) {
	a := assert.New(t)
	c, err := categoryService.GetCategoryBySlug(category.Slug)
	a.NoError(err)
	a.NotNil(c)
}

func TestCategoryService_GetCategoryByID(t *testing.T) {
	ass := assert.New(t)
	cat, err := categoryService.GetCategoryByID(1)
	ass.NoError(err)
	ass.NotNil(cat)
}

func TestCategoryService_GetCategoriesByQuery(t *testing.T) {
	ass := assert.New(t)
	cats, err := categoryService.GetCategoriesByQuery(repository.NewDefaultQuery())
	ass.NoError(err)
	ass.Len(cats, 10)
}