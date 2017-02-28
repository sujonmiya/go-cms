package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
	"models/capabilities"
	"models/roles"
)

func TestUploadsDir(t *testing.T) {
	a := assert.New(t)
	absolute, relative := UploadsDir()
	log.Printf("absolutePath: %s", absolute)
	log.Printf("relativePath: %s", relative)
	//a.Equal(`C:\Users\Sujon Miya\Documents\Projects\Contetto\cms\uploads\2016\11\27`, absolute)
	a.Equal("/uploads/2016/11/29", relative)
}

func TestAcl(t *testing.T) {
	ass := assert.New(t)
	ass.NotNil(Acl)
	ass.NotContains(Acl[roles.Subscriber], capabilities.UpdatePage)
	ass.Contains(Acl[roles.Editor], capabilities.UpdatePage)
	ass.NotContains(Acl[roles.Editor], capabilities.CreatePage)
	ass.Contains(Acl[roles.Author], capabilities.CreatePage)
	ass.NotContains(Acl[roles.Author], capabilities.CreateUser)
	ass.Contains(Acl[roles.Administrator], capabilities.DeleteUser)
}