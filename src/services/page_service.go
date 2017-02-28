package services

import (
	"log"
	"strings"

	"github.com/gosimple/slug"
	"github.com/microcosm-cc/bluemonday"
	"models"
	"models/status"
	"repository"
	"utils"
	"fmt"
	"controllers/viewmodels"
	"github.com/dustin/go-humanize"
	"strconv"
)

type PageService struct {
	repo *repository.Repository
}

func NewPageService() *PageService {
	return &PageService{repo: repository.NewRepo()}
}

func (s *PageService) SavePage(page *models.NewPage) (*models.Page, error) {
	title := strings.TrimSpace(page.Title)

	p := models.Page{
		Title:      title,
		Slug:       slug.Make(title),
		Content:    bluemonday.UGCPolicy().Sanitize(page.Content),
		Template:   strings.TrimSpace(page.Template),
		Status:     page.Status.String(),
		AuthorID:     utils.ToUInt32(page.Author.ID),
		LastEditorID:     utils.ToUInt32(page.Author.ID),
	}

	if err := s.repo.Save(&p); err != nil {
		log.Printf("Error creating Page: %v", err)
		return nil, err
	}

	return &p, nil
}

func (s *PageService) GetPageBySlug(slug string) (*models.Page, error) {
	filter := &models.Page{
		Slug:slug,
		Status:status.Draft.String(),
	}
	var page models.Page
	if err := s.repo.FindOne(filter, &page); err != nil {
		log.Printf("Error finding pages for slug %s : %v", slug, err)
		return nil, err
	}

	return &page, nil
}

func (s *PageService) GetPage(slug string) (*viewmodels.Page, error) {
	sql := `SELECT
              p.id                                        AS page_id,
			  p.title,
			  p.slug,
			  p.content,
			  p.template,
			  p.status,
			  p.created_at,
			  p.updated_at,
			  concat_ws(' ', u.first_name, u.last_name)   AS author_name,
			  concat_ws(' ', u2.first_name, u2.last_name) AS last_editor_name
			FROM pages AS p
			  INNER JOIN users AS u
			    ON p.author_id = u.id
			  INNER JOIN users AS u2
			    ON p.last_editor_id = u2.id
			WHERE p.slug = ?
			AND p.deleted_at IS NULL
			ORDER BY p.id ASC LIMIT 1`

	var result repository.PageResult
	err := s.repo.DB().Raw(sql, slug).Scan(&result).Error
	if err != nil {
		log.Printf("Error finding Page %s: %v", slug, err)
		return nil, err
	}

	page := &viewmodels.Page{}
	page.ID = strconv.Itoa(int(result.PageID))
	page.Title = result.Title
	page.Slug = result.Slug
	page.Content = result.Content
	page.Template = result.Template
	page.Status = result.Status
	page.AuthorName = result.AuthorName
	page.LastEditorName = result.LastEditorName
	page.CreatedAt = humanize.Time(result.CreatedAt)
	page.UpdatedAt = humanize.Time(result.UpdatedAt)

	return page, nil
}

func (s *PageService) GetPageByID(id uint32) (*viewmodels.Page, error) {
	sql := `SELECT
              p.id                                        AS page_id,
			  p.title,
			  p.slug,
			  p.content,
			  p.template,
			  p.status,
			  p.created_at,
			  p.updated_at,
			  concat_ws(' ', u.first_name, u.last_name)   AS author_name,
			  concat_ws(' ', u2.first_name, u2.last_name) AS last_editor_name
			FROM pages AS p
			  INNER JOIN users AS u
			    ON p.author_id = u.id
			  INNER JOIN users AS u2
			    ON p.last_editor_id = u2.id
			WHERE p.id = ?
			AND p.deleted_at IS NULL
			ORDER BY p.id ASC LIMIT 1`

	var result repository.PageResult
	err := s.repo.DB().Raw(sql, id).Scan(&result).Error
	if err != nil {
		log.Printf("Error finding Page by ID %s: %v", id, err)
		return nil, err
	}

	page := &viewmodels.Page{}
	page.ID = strconv.Itoa(int(result.PageID))
	page.Title = result.Title
	page.Slug = result.Slug
	page.Content = result.Content
	page.Template = result.Template
	page.Status = result.Status
	page.AuthorName = result.AuthorName
	page.LastEditorName = result.LastEditorName
	page.CreatedAt = humanize.Time(result.CreatedAt)
	page.UpdatedAt = humanize.Time(result.UpdatedAt)

	return page, nil
}

func (s *PageService) GetPages(query models.Query) ([]*models.Page, error) {
	var pages []*models.Page
	filter := &models.Page{
		Status:status.Draft.String(),
	}
	if err := s.repo.FindByQueryAndFilter(query, filter, &pages); err != nil {
		log.Printf("Error finding pages with query %+v: %v", query, err)
		return nil, err
	}

	return pages, nil
}

func (s *PageService) GetPagesByQuery(query models.Query) ([]*viewmodels.Page, error) {
	sql := `SELECT
              p.id                                        AS page_id,
			  p.title,
			  p.slug,
			  p.content,
			  p.template,
			  p.status,
			  p.created_at,
			  p.updated_at,
			  concat_ws(' ', u.first_name, u.last_name)   AS author_name,
			  concat_ws(' ', u2.first_name, u2.last_name) AS last_editor_name
			FROM pages AS p
			  INNER JOIN users AS u
			    ON p.author_id = u.id
			  INNER JOIN users AS u2
			    ON p.last_editor_id = u2.id
			WHERE p.deleted_at IS NULL`

	var result []repository.PageResult
	err := s.repo.DB().Raw(sql).
		Limit(query.Total).
		Offset(query.Offset).
		Order(fmt.Sprintf("p.%s", query.Sort)).
		Scan(&result).Error
	if err != nil {
		log.Printf("Error finding Pages with Query %#v: %v", query, err)
		return nil, err
	}

	pages := []*viewmodels.Page{}
	for _, r := range result {
		page := &viewmodels.Page{}
		page.ID = strconv.Itoa(int(r.PageID))
		page.Title = r.Title
		page.Slug = r.Slug
		page.Content = r.Content
		page.Template = r.Template
		page.Status = r.Status
		page.AuthorName = r.AuthorName
		page.LastEditorName = r.LastEditorName
		page.CreatedAt = humanize.Time(r.CreatedAt)
		page.UpdatedAt = humanize.Time(r.UpdatedAt)

		pages = append(pages, page)
	}

	return pages, nil
}