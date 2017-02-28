package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"utils/seed"
)

var taxService *TaxonomyService

func init() {
	_service := NewTaxonomyService()
	taxService = _service
}

func TestNewTaxonomyService(t *testing.T) {
	a := assert.New(t)
	a.NotNil(taxService)
}

func TestTaxonomyService_SaveTaxonomy(t *testing.T) {
	ass := assert.New(t)
	tax := seed.NewTaxonomy()
	tax.Parent = 4
	taxonomy, err := taxService.SaveTaxonomy(tax)
	ass.NoError(err)
	ass.NotNil(taxonomy)
}
