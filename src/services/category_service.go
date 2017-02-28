package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/mgo.v2/bson"
	"models"
	"utils"
	"repository"
	"controllers/viewmodels"
	"github.com/dustin/go-humanize"
	"strconv"
)

type CategoryService struct {
	repo *repository.Repository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{repo: repository.NewRepo()}
}

func (s *CategoryService) SaveCategory(category *models.NewCategory) (*models.Category, error) {
	name := strings.TrimSpace(category.Name)
	cat := &models.Category{
		Name:        name,
		Slug:        slug.Make(name),
		Description: bluemonday.UGCPolicy().Sanitize(category.Description),
		AuthorID:     utils.ToUInt32(category.Author.ID),
		LastEditorID:     utils.ToUInt32(category.Author.ID),
	}

	if category.Parent != 0 {
		cat.ParentID = utils.ToUInt32(category.Parent)
	}

	if err := s.repo.Save(cat); err != nil {
		log.Printf("Error creating Category: %v", err)
		return nil, err
	}

	return cat, nil
}

func (s *CategoryService) UpdateCategory(category *models.UpdateCategory) error {
	category.UpdatedAt = time.Now()
	category.Description = bluemonday.UGCPolicy().Sanitize(category.Description)
	/*if err := s.repo.Update(category.Id, category); err != nil {
		log.Printf("error updating category: %v - %#v", err, category)
		return err
	}*/

	return nil
}

func (s *CategoryService) UpdateCategories(categories []*models.UpdateCategory) error {
	var errs []string
	for _, category := range categories {
		if err := s.UpdateCategory(category); err != nil {
			errs = append(errs, err.Error())
			return err
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ", "))
	}

	return nil
}

/*func (s *CategoryService) GetCategoryByID(id uint32) (*models.Category, error) {
	var category models.Category
	if err := s.repo.FindByID(id, &category); err != nil {
		log.Printf("Error finding category ID#%s: %v", id, err)
		return nil, err
	}

	return &category, nil
}*/

func (s *CategoryService) GetCategoryBySlug(slug string) (*models.Category, error) {
	filter := &models.Category{Slug:slug}
	var category models.Category
	if err := s.repo.FindOne(filter, &category); err != nil {
		log.Printf("error finding category for slug: %s %v", slug, err)
		return nil, err
	}

	return &category, nil
}

func (s *CategoryService) GetCategoryByID(id uint32) (*viewmodels.Category, error) {
	sql := `SELECT
		  c.id                                      AS category_id,
		  c.name,
		  c.slug,
		  c.description,
		  (SELECT
		     GROUP_CONCAT(DISTINCT CONCAT_WS(':', p.id, p.name, p.slug) ORDER BY p.name)
		  FROM categories AS p
		  WHERE p.id = c.parent_id) 			AS parent,
		  (SELECT count(*)
		  FROM article_categories AS ac
		  WHERE ac.category_id = c.id) 			AS num_articles,
		  c.created_at,
		  c.updated_at,
		  u.id                                       AS author_id,
		  CONCAT_WS(' ', u.first_name, u.last_name)  AS author_name,
		  CONCAT_WS(' ', u.first_name, u2.last_name) AS last_editor_name
		FROM categories AS c
		  INNER JOIN users AS u
		    ON c.author_id = u.id
		  INNER JOIN users AS u2
		    ON c.last_editor_id = u2.id
		WHERE c.deleted_at IS NULL
		AND c.id = ? LIMIT 1`

	var result repository.CategoryResult
	err := s.repo.DB().Raw(sql, id).Scan(&result).Error
	if err != nil {
		log.Printf("Error finding Category by ID %+v: %v", id, err)
		return nil, err
	}

	category := &viewmodels.Category{}
	category.ID = strconv.Itoa(int(result.CategoryID))
	category.Name = result.Name
	category.Slug = result.Slug
	category.Description = result.Description
	category.Parent = toIDNameSlug(result.Parent)
	category.NumArticles = result.NumArticles
	category.AuthorName = result.AuthorName
	category.LastEditorName = result.LastEditorName
	category.CreatedAt = humanize.Time(result.CreatedAt)
	category.UpdatedAt = humanize.Time(result.UpdatedAt)

	return category, nil
}

func (s *CategoryService) GetArticlesFromCategory(slug string) ([]*models.Article, error) {
	filter := &models.Article{Slug:slug}
	var articles []*models.Article
	if err := s.repo.FindByQueryAndFilter(repository.NewDefaultQuery(), filter, &articles); err != nil {
		log.Printf("error finding articles for category slug: %s %v", slug, err)
		return nil, err
	}

	return articles, nil
}

func (s *CategoryService) GetCategories() ([]*viewmodels.Category, error) {
	var categories []*viewmodels.Category
	/*if err := s.repo.FindByObjectIds(ids, &categories); err != nil {
		log.Printf("error finding categories %v: %v", ids, err)
		return nil, err
	}*/

	return categories, nil
}

func (s *CategoryService) GetCategoriesByIDs(ids []bson.ObjectId) ([]*models.Category, error) {
	var categories []*models.Category
	/*if err := s.repo.FindByObjectIds(ids, &categories); err != nil {
		log.Printf("error finding categories %v: %v", ids, err)
		return nil, err
	}*/

	return categories, nil
}

func (s *CategoryService) GetCategoriesByQuery(query models.Query) ([]*viewmodels.Category, error) {
	sql := `SELECT
			  c.id                                      AS category_id,
			  c.name,
			  c.slug,
			  c.description,
			  (SELECT
			     GROUP_CONCAT(DISTINCT CONCAT_WS(':', p.id, p.name, p.slug) ORDER BY p.name)
			  FROM categories AS p
			  WHERE p.id = c.parent_id) AS parent,
			  (SELECT count(*)
			  FROM article_categories AS ac
			  WHERE ac.category_id = c.id) AS num_articles,
              c.created_at,
			  c.updated_at,
			  u.id                                       AS author_id,
			  CONCAT_WS(' ', u.first_name, u.last_name)  AS author_name,
			  CONCAT_WS(' ', u.first_name, u2.last_name) AS last_editor_name
			FROM categories AS c
			  INNER JOIN users AS u
			    ON c.author_id = u.id
			  INNER JOIN users AS u2
			    ON c.last_editor_id = u2.id
			WHERE c.deleted_at IS NULL`

	var result []*repository.CategoryResult
	err := s.repo.DB().Raw(sql).
		Limit(query.Total).
		Offset(query.Offset).
		Order(fmt.Sprintf("c.%s", query.Sort)).
		Scan(&result).Error
	if err != nil {
		log.Printf("Error finding Categories by Query %+v: %v", query, err)
		return nil, err
	}

	categories := []*viewmodels.Category{}
	for _, c := range result {
		cat := &viewmodels.Category{}
		cat.ID = strconv.Itoa(int(c.CategoryID))
		cat.Name = c.Name
		cat.Slug = c.Slug
		cat.Description = c.Description
		cat.Parent = toIDNameSlug(c.Parent)
		cat.NumArticles = c.NumArticles
		cat.AuthorName = c.AuthorName
		cat.LastEditorName = c.LastEditorName
		cat.CreatedAt = humanize.Time(c.CreatedAt)
		cat.UpdatedAt = humanize.Time(c.UpdatedAt)

		categories = append(categories, cat)
	}

	return categories, nil
}

/*
func (s *CategoryService) GetAllCategories() ([]*models.Category, error) {
	var categories []*models.Category
	if err := s.repo.FindAll(&categories); err != nil {
		log.Printf("error finding categories: %v", err)
		return nil, err
	}

	return categories, nil
}

func (s *CategoryService) DeleteCategory(id string) error {
	// Rx: implement tx
	if !bson.IsObjectIdHex(id) {
		return ErrInvalidArticleId
	}

	return s.repo.Delete(bson.ObjectIdHex(id))
}*/
