package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"utils/seed"
	"repository"
	"models/status"
)

var (
	articleService *ArticleService
	articleId string
)

func init() {
	_service := NewArticleService()
	articleService = _service
}

func TestArticleService(t *testing.T) {
	a := assert.New(t)
	a.NotNil(articleService)
}

func TestArticleService_SaveArticle(t *testing.T) {
	ass := assert.New(t)
	art := seed.NewArticle()
	art.Categories = []uint32{1, 2, 3, 5}
	art.Taxonomies = []uint32{4, 5, 6}
	article, err := articleService.SaveArticle(art)
	ass.NoError(err)
	ass.NotNil(article)
}

func TestArticleService_GetArticleBySlug(t *testing.T) {
	ass := assert.New(t)
	slug := "fugit-et-magnam-ut"
	article, err := articleService.GetArticleBySlug(slug)
	ass.NoError(err)
	ass.NotNil(article)
}

func TestArticleService_GetArticleByID(t *testing.T) {
	ass := assert.New(t)
	article, err := articleService.GetArticleByID(1)
	ass.NoError(err)
	ass.NotNil(article)
}

func TestArticleService_GetArticlesByQuery(t *testing.T) {
	ass := assert.New(t)
	articles, err := articleService.GetArticlesByQuery(repository.NewDefaultQuery())
	ass.NoError(err)
	ass.Len(articles, 6)
}

func TestArticleService_GetArticlesByFilter(t *testing.T) {
	ass := assert.New(t)
	filter := repository.Filter{
		Author: []uint32{2},
		//Editor: []uint32{1,2},
		//Category: []uint32{1,2,3,4,5},
		//Taxonomy: []uint32{1,2,3,4,5},
		Status:status.Draft.String(),
		Total:50,
		//Sort:[]string{"id"},
	}
	articles, err := articleService.GetArticlesByFilter(filter)
	ass.NoError(err)
	ass.Len(articles, 5)
}

func TestArticleService_ToNameSlug(t *testing.T) {
	ass := assert.New(t)
	result := toIDNameSlug("culpa:culpa,non:non,quae:quae")
	ass.Len(result, 3)
}

/*func TestGetArticle(t *testing.T) {
	assert := assert.New(t)

	article, _ := articleService.GetArticle(articleId)
	assert.NotNil(article)
}

func TestGetArticleNotFoundError(t *testing.T) {
	assert := assert.New(t)

	_, err := articleService.GetArticle(bson.NewObjectId().Hex())
	assert.Equal(mgo.ErrNotFound, err)
}

func TestGetArticleForSlug(t *testing.T) {
	assert := assert.New(t)

	article, _ := articleService.GetArticleForSlug("accusantium-odit-quis-ut")
	assert.NotNil(article)
}

func TestGetArticlesFromAuthor(t *testing.T) {
	assert := assert.New(t)

	article, _ := articleService.GetArticlesFromAuthor(bson.ObjectIdHex("56c3047e6cbee702b8ffa08b"))
	assert.NotNil(article)
}

func TestGetArticlesFromCategory(t *testing.T) {
	assert := assert.New(t)

	articles, _ := articleService.GetArticlesFromCategory(bson.ObjectIdHex("56c3047e6cbee702b8ffa08b"))
	assert.Empty(articles)
}

func TestGetArticlesFromTag(t *testing.T) {
	assert := assert.New(t)

	articles, _ := articleService.GetArticlesFromTag(bson.ObjectIdHex("56c3047e6cbee702b8ffa08b"))
	assert.Empty(articles)
}

func TestDeleteArticle(t *testing.T) {
	assert := assert.New(t)

	err := articleService.DeleteArticle(articleId)
	assert.NoError(err)
}*/
