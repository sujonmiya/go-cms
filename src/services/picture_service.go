package services

import (
	"models"
	"log"
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	_ "image/gif"
	"path/filepath"
	"config"
	"os"
	"fmt"
	"strings"
	"repository"
	"controllers/viewmodels"
	"strconv"
	"github.com/dustin/go-humanize"
)

type PictureService struct {
	repo *repository.Repository
	*FileService
}

func NewPictureService() *PictureService {
	return &PictureService{
		repo:repository.NewRepo(),
		FileService:NewFileService(),
	}
}

func (s *PictureService) GetPictures() ([]*models.Picture, error) {
	var pics []*models.Picture
	if err := s.repo.Find(&pics); err != nil {
		return nil, err
	}

	return pics, nil
}

func (s *PictureService) GetPicturesByQuery(query models.Query) ([]*viewmodels.Picture, error) {
	sql := `SELECT
		  p.id                                       AS picture_id
		  ,p.name
		  ,p.width
		  ,p.height
		  ,p.size
		  ,p.caption
		  ,p.alt_text
		  ,p.mime_type
		  ,p.url
		  ,u.id                                      AS uploader_id
		  ,CONCAT_WS(' ', u.first_name, u.last_name) AS uploader_name
		  ,u.created_at
		  ,u.updated_at
		FROM pictures AS p
		  LEFT JOIN users AS u
		    ON p.author_id = u.id
		WHERE p.deleted_at IS NULL`
	var result []repository.PictureResult
	err := s.repo.DB().Raw(sql).
		Limit(query.Total).
		Offset(query.Offset).
		Order(fmt.Sprintf("p.%s", query.Sort)).
		Scan(&result).Error
	if err != nil {
		log.Printf("Error finding Pictures by Query %+v: %v", query, err)
		return nil, err
	}

	pics := []*viewmodels.Picture{}
	for _, p := range result {
		pic := &viewmodels.Picture{}
		pic.ID = strconv.Itoa(int(p.PictureID))
		pic.Name = p.Name
		pic.Width = p.Width
		pic.Height = p.Height
		pic.Size = humanize.Bytes(uint64(p.Size))
		pic.Caption = p.Caption
		pic.Description = p.AltText
		pic.MimeType = p.MimeType
		pic.Url = p.Url
		pic.Uploader.ID = strconv.Itoa(int(p.UploaderID))
		pic.Uploader.Name = p.UploaderName
		pic.UploadedAt = humanize.Time(p.CreatedAt)
		pic.UpdatedAt = humanize.Time(p.UpdatedAt)

		pics = append(pics, pic)
	}

	return pics, nil
}

func (s *PictureService) CreateAndSavePicture(picture *models.NewPicture) (*models.Picture, error) {
	p := &models.Picture{
		Name: picture.Name,
		Caption: picture.Caption,
		AltText: picture.AltText,
		MimeType: picture.MimeType,
		Size:len(picture.Data),
	}

	conf, _, err := image.DecodeConfig(bytes.NewReader(picture.Data))
	if err != nil {
		log.Printf("Error decoding the image: %v", err)
		return nil, err
	}

	p.Width = conf.Width
	p.Height = conf.Height
	hash, err := sha256Checksum(picture.Data);
	pictureName := generatePictureName(picture.Name, hash, conf.Width, conf.Height)
	if err != nil {
		log.Printf("Error calculating checksum of picture file %s: %v", pictureName, err)
		return nil, err
	}

	absolutePath, relativePath := config.UploadsDir()
	filename := filepath.Join(absolutePath, pictureName)
	if s.FileExists(filename) {
		log.Printf("Error picture file %s already exist\n", filename)
		return nil, os.ErrExist
	}

	if err := s.WriteFile(filename, picture.Data); err != nil {
		return nil, err
	}

	p.Url = fmt.Sprintf("%s/%s", relativePath, pictureName)
	if err := s.repo.Save(p); err != nil {
		log.Printf("Error saving Picture: %v", err)
		return nil, err
	}

	return p, nil
}

func (s *PictureService) DeletePicture(picture *models.Picture) error {
	tx := s.repo.DB().Begin()
	filename := filepath.Join(config.BasePath(), picture.Url)
	if err := s.RemoveFile(filename); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(picture).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting Picture: %v - \n%#v\n", err, picture)

		return err
	}
	tx.Commit()

	return nil
}

func generatePictureName(s, hash string, width, height int) string {
	ext := filepath.Ext(s)
	basename := strings.TrimSuffix(s, ext)
	return fmt.Sprintf("%s_%s_%dx%d%s", basename, hash, width, height, ext)
}