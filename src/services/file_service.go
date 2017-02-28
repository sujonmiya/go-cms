package services

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"crypto/sha256"
	"crypto/hmac"
	"encoding/base64"
)

const (
	fileMode = os.FileMode(0664)
	macKey = `:2#-^/Z6)8)%p/56+_.YS]};"9;+4:>c`
)

type FileService struct {}

func NewFileService() *FileService {
	return &FileService{}
}

func (s *FileService) FileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func (s *FileService) WriteFile(filename string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), fileMode); err != nil {
		log.Printf("Error creating dir %s: %v\n", filename, err)
		return err
	}

	if err := ioutil.WriteFile(filename, data, fileMode); err != nil {
		log.Printf("Error writing file to %s: %v\n", filename, err)
		return err
	}

	return nil
}

func (s *FileService) RemoveFile(filename string) error {
	if err := os.Remove(filename); err != nil {
		log.Printf("Error removing file %s: %v\n", filename, err)
		return err
	}

	return nil
}

func sha256Checksum(data []byte) (string, error) {
	mac := hmac.New(sha256.New, []byte(macKey))
	_, err := mac.Write(data)
	if err != nil {
		return "", err
	}

	hashed := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	return hashed, nil
}