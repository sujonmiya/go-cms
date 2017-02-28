package services

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"time"
	"strconv"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

var (
	fileService *FileService
	filename = `C:\Users\Sujon Miya\Documents\Projects\Contetto\cms\uploads\2016\11\25\icon-user-default.png`
)

func init() {
	_service := NewFileService()
	fileService = _service
}

func TestNewFileService(t *testing.T) {
	ass := assert.New(t)
	ass.NotNil(fileService)
}

func TestFileService_FileExists(t *testing.T) {
	a := assert.New(t)
	ok := fileService.FileExists(filename)
	a.True(ok)
}

func TestFileService_WriteFile(t *testing.T) {
	ass := assert.New(t)
	data, _ := ioutil.ReadFile(filename)
	day := strconv.Itoa(time.Now().Day())
	err := fileService.WriteFile(strings.Replace(filename, "25", day, 1), data)
	ass.NoError(err)
}

func TestFileService_RemoveFile(t *testing.T) {
	ass := assert.New(t)
	day := strconv.Itoa(time.Now().Day())
	err := fileService.RemoveFile(strings.Replace(filename, "25", day, 1))
	ass.NoError(err)
}

func TestFileService_Sha256Checksum(t *testing.T) {
	ass := assert.New(t)
	message := []byte("Hi, there!")
	hashed, err := sha256Checksum(message)
	ass.NoError(err)

	mac := hmac.New(sha256.New, []byte("wrongkey"))
	mac.Write(message)
	toCheck := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	ass.NotEqual(hashed, toCheck)
	decoded, _ := base64.URLEncoding.DecodeString(toCheck)
	ass.NotEqual(string(message), string(decoded))
}
