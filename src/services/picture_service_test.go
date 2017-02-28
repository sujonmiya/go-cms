package services

import (
	"models"
	"testing"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"io/ioutil"
	"github.com/icrowley/fake"
	"net/http"
	"log"
	"strings"
	"repository"
)

var (
	pictureService *PictureService
	picture *models.Picture
)

func init() {
	_service := NewPictureService()
	pictureService = _service
}

func TestNewPictureService(t *testing.T) {
	a := assert.New(t)
	a.NotNil(pageService)
}

func TestPictureService_GetPictures(t *testing.T) {
	a := assert.New(t)
	pics, err := pictureService.GetPictures()
	a.NoError(err)
	a.NotEmpty(pics)
}

func TestPictureService_SavePicture(t *testing.T) {
	a := assert.New(t)
	filename := `C:\Users\Sujon Miya\Documents\Projects\Contetto\cms\uploads\2016\11\25\icon-user-default.png`
	data, _ := ioutil.ReadFile(filename)
	newPic := &models.NewPicture{
		Name:filepath.Base(filename),
		Data:data,
		AltText:fake.Words(),
		Caption:fake.Words(),
		MimeType: http.DetectContentType(data),
	}

	_picture, err := pictureService.CreateAndSavePicture(newPic)
	a.NoError(err)
	a.NotNil(_picture)
	picture = _picture
}

func TestPictureService_SavePictureUnknownFormat(t *testing.T) {
	a := assert.New(t)
	filename := `C:\Users\Sujon Miya\Documents\Projects\Contetto\cms\templates\assets\css\style.css`
	data, _ := ioutil.ReadFile(filename)
	newPic := &models.NewPicture{
		Data:data,
	}

	_, err := pictureService.CreateAndSavePicture(newPic)
	a.Error(err)
}

func TestPictureService_DeletePicture(t *testing.T) {
	a := assert.New(t)
	pic := &models.Picture{Model:models.Model{ID:1}}
	err := pictureService.DeletePicture(pic)
	a.NoError(err)
}

func TestGeneratePictureName(t *testing.T) {
	ass := assert.New(t)
	s := "icon-user-default.png"
	ext := filepath.Ext(s)
	basename := strings.TrimSuffix(s, ext)
	log.Printf("basename: %s", basename)
	log.Printf("extension: %s", ext)
	ass.Equal("icon-user-default", basename)
	ass.Equal(".png", ext)
}

func TestPictureService_GetPicturesByQuery(t *testing.T) {
	ass := assert.New(t)
	pics, err := pictureService.GetPicturesByQuery(repository.NewDefaultQuery())
	ass.NoError(err)
	ass.Len(pics, 9)
}