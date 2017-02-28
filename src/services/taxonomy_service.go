package services

import (
	"github.com/gosimple/slug"
	"github.com/microcosm-cc/bluemonday"
	"models"
	"log"
	"strings"
	"utils"
	"repository"
)

type TaxonomyService struct {
	repo *repository.Repository
}

func NewTaxonomyService() *TaxonomyService {
	return &TaxonomyService{repo: repository.NewRepo()}
}

func (s *TaxonomyService) SaveTaxonomy(taxonomy *models.NewTaxonomy) (*models.Taxonomy, error) {
	name := strings.TrimSpace(taxonomy.Name)
	tax := &models.Taxonomy{
		Name:        name,
		Slug:        slug.Make(name),
		Description: bluemonday.UGCPolicy().Sanitize(taxonomy.Description),
		AuthorID:     utils.ToUInt32(taxonomy.Author.ID),
		LastEditorID:     utils.ToUInt32(taxonomy.Author.ID),
	}

	if taxonomy.Parent != 0 {
		tax.ParentID = utils.ToUInt32(taxonomy.Parent)
	}

	if err := s.repo.Save(tax); err != nil {
		log.Printf("Error creating Taxonomy: %v", err)
		return nil, err
	}

	return tax, nil
}

func (s *TaxonomyService) GetTaxonomyByID(id uint32) (*models.Taxonomy, error) {
	var tax models.Taxonomy
	if err := s.repo.FindByID(id, &tax); err != nil {
		log.Printf("Error finding taxonomy ID#%s: %v", id, err)
		return nil, err
	}

	return &tax, nil
}
/*
func (s *TagService) GetTags(ids []bson.ObjectId) ([]*models.Tag, error) {
	var tags []*models.Tag
	if err := s.repo.FindByObjectIds(ids, &tags); err != nil {
		log.Printf("error finding tags #", ids, " : ", err)
		return nil, err
	}

	return tags, nil
}

func (s *TagService) GetAllTags() ([]*models.Tag, error) {
	var tags []*models.Tag
	if err := s.repo.FindAll(&tags); err != nil {
		log.Printf("error finding tags: ", err)
		return nil, err
	}

	return tags, nil
}

func (s *TagService) UpdateTag(tag *models.Tag) error {
	tag.UpdatedAt = time.Now()
	tag.Description = bluemonday.UGCPolicy().Sanitize(tag.Description)
	if err := s.repo.Update(tag.ID, tag); err != nil {
		log.Printf("error updating tag: ", err, " - ", tag)
		return err
	}

	return nil
}*/
