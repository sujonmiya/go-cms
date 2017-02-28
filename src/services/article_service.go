package services

import (
	"log"
	"models"

	"github.com/gosimple/slug"
	"github.com/microcosm-cc/bluemonday"
	"models/status"
	"strings"
	"config"
	"path/filepath"
	"utils"
	"repository"
	"controllers/viewmodels"
	"fmt"
	"github.com/dustin/go-humanize"
	"strconv"
)

type ArticleService struct {
	repo *repository.Repository
	cs   *CategoryService
	ts   *TaxonomyService
	ps   *PictureService
}

func NewArticleService() *ArticleService {
	return &ArticleService{
		repo: repository.NewRepo(),
		cs: NewCategoryService(),
		ts: NewTaxonomyService(),
		ps: NewPictureService(),
	}
}

func (s *ArticleService) SaveArticle(article *models.NewArticle) (*models.Article, error) {
	title := strings.TrimSpace(article.Title)
	art := &models.Article{
		Title:      title,
		Slug:       slug.Make(title),
		Content:    bluemonday.UGCPolicy().Sanitize(article.Content),
		Excerpt: article.Excerpt,
		Status:     article.Status.String(),
		AuthorID:     utils.ToUInt32(article.Author.ID),
		LastEditorID:     utils.ToUInt32(article.Author.ID),
	}

	tx := s.repo.DB().Begin()
	if article.HasFeaturedImage() {
		picture, err := s.ps.CreateAndSavePicture(article.FeaturedImage)
		if err != nil {
			tx.Rollback()
			log.Printf("Error saving article Picture: %v", err)
			return nil, err
		}

		art.Picture = picture
		art.Picture.AuthorID = utils.ToUInt32(article.Author.ID)
		art.Picture.LastEditorID = utils.ToUInt32(article.Author.ID)
	}

	if err := tx.Create(art).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating Article: %v", err)
		filename := filepath.Join(config.BasePath(), art.Picture.Url)
		if err := s.ps.RemoveFile(filename); err != nil {
			tx.Rollback()
			return nil, err
		}

		return nil, err
	}

	for _, id := range article.Categories {
		ac := &repository.ArticleCategory{
			ArticleID: utils.ToUInt32(art.ID),
			CategoryID: utils.ToUInt32(id)}
		if err := tx.Save(ac).Error; err != nil {
			tx.Rollback()
			log.Printf("Error saving ArticleCategory: %v", err)
			return nil, err
		}
	}

	for _, id := range article.Taxonomies {
		at := &repository.ArticleTaxonomy{
			ArticleID: utils.ToUInt32(art.ID),
			TaxonomyID: utils.ToUInt32(id)}
		if err := tx.Create(at).Error; err != nil {
			tx.Rollback()
			log.Printf("Error saving ArticleTaxonomy: %v", err)
			return nil, err
		}
	}

	tx.Commit()
	return art, nil
}

func (s *ArticleService) GetArticleForSlug(slug string) (*models.Article, error) {
	filter := &models.Article{
		Slug:slug,
		Status:status.Draft.String(),
	}

	var article models.Article
	if err := s.repo.FindOne(filter, &article); err != nil {
		log.Printf("Error finding article for slug %s: %v", slug, err)
		return nil, err
	}

	return &article, nil
}

func (s *ArticleService) GetArticles(query models.Query) ([]*models.Article, error) {
	var articles []*models.Article
	filter := &models.Article{
		Status:status.Draft.String(),
	}

	if err := s.repo.FindByQueryAndFilter(query, filter, &articles); err != nil {
		log.Printf("Error finding articles with Query %+v: %v", query, err)
		return nil, err
	}

	return articles, nil
}

func (s *ArticleService) GetArticlesByQuery(query models.Query) ([]*viewmodels.Article, error) {
	var result []repository.ArticleResult
	err := s.repo.DB().Table(repository.TableViewArticles).
		Limit(query.Total).
		Offset(query.Offset).
		Order(fmt.Sprintf("a.%s", query.Sort)).
		Find(&result).Error
	if err != nil {
		log.Printf("Error finding Articles by Query %+v: %v", query, err)
		return nil, err
	}

	articles := []*viewmodels.Article{}
	for _, a := range result {
		art := &viewmodels.Article{}
		art.ID = strconv.Itoa(int(a.ArticleID))
		art.Title = a.Title
		art.Slug = a.Slug
		art.Content = a.Content
		art.Excerpt = a.Excerpt

		art.FeaturedImage.ID = strconv.Itoa(int(a.FeaturedImageID))
		art.FeaturedImage.Caption = a.FeaturedImageCaption
		art.FeaturedImage.AltText = a.FeaturedImageCaption
		art.FeaturedImage.Url = a.FeaturedImageUrl

		art.Author.ID = strconv.Itoa(int(a.AuthorID))
		art.Author.FullName = a.AuthorName
		art.Author.NickName = a.AuthorNickName
		art.Author.Website = a.AuthorWebsite
		art.Author.Biography = a.AuthorBio
		art.Author.ProfilePictureUrl = a.AuthorProfilePicUrl

		art.Editor.ID = strconv.Itoa(int(a.EditorID))
		art.Editor.FullName = a.EditorName
		art.Editor.NickName = a.EditorNickName

		art.Categories = toIDNameSlugs(a.Categories)
		art.Taxonomies = toIDNameSlugs(a.Taxonomies)
		art.Status = a.Status
		art.CreatedAt = humanize.Time(a.CreatedAt)
		art.UpdatedAt = humanize.Time(a.UpdatedAt)

		articles = append(articles, art)
	}

	return articles, nil
}

func (s *ArticleService) GetArticlesByFilter(f repository.Filter) ([]*viewmodels.Article, error) {
	t := s.repo.DB().
		Table(repository.TableViewArticles).
		Limit(f.Total).
		Offset(f.Offset)
	if len(f.Author) > 0 {
		t = t.Where("`vw_articles`.`author_id` IN (?)", f.Author)
	}

	if len(f.Editor) > 0 {
		t = t.Where("`vw_articles`.`editor_id` IN (?)", f.Editor)
	}

	if len(f.Category) > 0 {
		t = t.Where("`vw_articles`.`categories` IN (?)", f.Category)
	}

	if len(f.Taxonomy) > 0 {
		t = t.Where("`vw_articles`.`taxonomies` IN (?)", f.Taxonomy)
	}

	if f.Status != "" {
		t = t.Where("`vw_articles`.`status` = ?", f.Status)
	}

	var result []repository.ArticleResult
	if err := t.Find(&result).Error; err != nil {
		log.Printf("Error finding Articles by Filter %+v: %v", f, err)
		return nil, err
	}

	articles := []*viewmodels.Article{}
	for _, a := range result {
		art := &viewmodels.Article{}
		art.ID = strconv.Itoa(int(a.ArticleID))
		art.Title = a.Title
		art.Slug = a.Slug
		art.Content = a.Content
		art.Excerpt = a.Excerpt

		art.FeaturedImage.ID = strconv.Itoa(int(a.FeaturedImageID))
		art.FeaturedImage.Caption = a.FeaturedImageCaption
		art.FeaturedImage.AltText = a.FeaturedImageCaption
		art.FeaturedImage.Url = a.FeaturedImageUrl

		art.Author.ID = strconv.Itoa(int(a.AuthorID))
		art.Author.FullName = a.AuthorName
		art.Author.NickName = a.AuthorNickName
		art.Author.Website = a.AuthorWebsite
		art.Author.Biography = a.AuthorBio
		art.Author.ProfilePictureUrl = a.AuthorProfilePicUrl

		art.Editor.ID = strconv.Itoa(int(a.EditorID))
		art.Editor.FullName = a.EditorName
		art.Editor.NickName = a.EditorNickName

		art.Categories = toIDNameSlugs(a.Categories)
		art.Taxonomies = toIDNameSlugs(a.Taxonomies)
		art.Status = a.Status
		art.CreatedAt = humanize.Time(a.CreatedAt)
		art.UpdatedAt = humanize.Time(a.UpdatedAt)

		articles = append(articles, art)
	}

	return articles, nil
}

func (s *ArticleService) GetArticleBySlug(slug string) (*viewmodels.Article, error) {
	var result repository.ArticleResult
	if err := s.repo.DB().
		Table(repository.TableViewArticles).
		Where("`vw_articles`.`slug` = ?", slug).
		First(&result).Error; err != nil {
		log.Printf("Error finding Article by slug %s: %v", slug, err)
		return nil, err
	}

	article := &viewmodels.Article{}
	article.ID = strconv.Itoa(int(result.ArticleID))
	article.Title = result.Title
	article.Slug = result.Slug
	article.Content = result.Content
	article.Excerpt = result.Excerpt

	article.FeaturedImage.ID = strconv.Itoa(int(result.FeaturedImageID))
	article.FeaturedImage.Caption = result.FeaturedImageCaption
	article.FeaturedImage.AltText = result.FeaturedImageCaption
	article.FeaturedImage.Url = result.FeaturedImageUrl

	article.Author.ID = strconv.Itoa(int(result.AuthorID))
	article.Author.FullName = result.AuthorName
	article.Author.NickName = result.AuthorNickName
	article.Author.Website = result.AuthorWebsite
	article.Author.Biography = result.AuthorBio
	article.Author.ProfilePictureUrl = result.AuthorProfilePicUrl

	article.Editor.ID = strconv.Itoa(int(result.EditorID))
	article.Editor.FullName = result.EditorName
	article.Editor.NickName = result.EditorNickName

	article.Categories = toIDNameSlugs(result.Categories)
	article.Taxonomies = toIDNameSlugs(result.Taxonomies)
	article.Status = result.Status
	article.CreatedAt = humanize.Time(result.CreatedAt)
	article.UpdatedAt = humanize.Time(result.UpdatedAt)

	return article, nil
}
func (s *ArticleService) GetArticleByID(id uint32) (*viewmodels.Article, error) {
	var result repository.ArticleResult
	if err := s.repo.DB().
		Table("vw_articles").
		Where("`vw_articles`.`article_id` = ?", id).
		First(&result).Error; err != nil {
		log.Printf("Error finding Article by ID %s: %v", id, err)
		return nil, err
	}

	article := &viewmodels.Article{}
	article.ID = strconv.Itoa(int(result.ArticleID))
	article.Title = result.Title
	article.Slug = result.Slug
	article.Content = result.Content
	article.Excerpt = result.Excerpt

	article.FeaturedImage.ID = strconv.Itoa(int(result.FeaturedImageID))
	article.FeaturedImage.Name = result.FeaturedImageName
	article.FeaturedImage.Caption = result.FeaturedImageCaption
	article.FeaturedImage.AltText = result.FeaturedImageDesc
	article.FeaturedImage.Width = result.FeaturedImageWidth
	article.FeaturedImage.Height = result.FeaturedImageHeight
	article.FeaturedImage.Size = humanize.Bytes(uint64(result.FeaturedImageSize))
	article.FeaturedImage.Url = result.FeaturedImageUrl

	article.Author.ID = strconv.Itoa(int(result.AuthorID))
	article.Author.FullName = result.AuthorName
	article.Author.NickName = result.AuthorNickName
	article.Author.Website = result.AuthorWebsite
	article.Author.Biography = result.AuthorBio
	article.Author.ProfilePictureUrl = result.AuthorProfilePicUrl

	article.Editor.ID = strconv.Itoa(int(result.EditorID))
	article.Editor.FullName = result.EditorName
	article.Editor.NickName = result.EditorNickName

	article.Categories = toIDNameSlugs(result.Categories)
	article.Taxonomies = toIDNameSlugs(result.Taxonomies)
	article.Status = result.Status
	article.CreatedAt = humanize.Time(result.CreatedAt)
	article.UpdatedAt = humanize.Time(result.UpdatedAt)

	return article, nil
}

func toIDNameSlug(s string) struct{ID, Name, Slug string} {
	result := struct{ID, Name, Slug string}{}
	pair := strings.Split(s, ":")
	if len(pair) == 3 {
		result.ID = pair[0]
		result.Name = pair[1]
		result.Slug = pair[2]
	}

	return result
}

func toIDNameSlugs(s string) []struct{ ID, Name, Slug string} {
	result := []struct{ID, Name, Slug string}{}
	for _, s := range strings.Split(s, ",") {
		result = append(result, toIDNameSlug(s))
	}

	return result
}
