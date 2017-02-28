package models

import (
	"testing"
	"time"

	"github.com/enodata/faker"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestHasCapability(t *testing.T) {
	assert := assert.New(t)

	user := User{
		Id:        bson.NewObjectId(),
		FirstName: faker.Name().FirstName(),
		LastName:  faker.Name().LastName(),
		Email:     faker.Internet().Email(),
		Password:  []byte(faker.Internet().Password(6, 18)),
		Role:      role.Author,
		//Website:faker.Internet().DomainName(),
		ShortBio:       faker.Lorem().Paragraph(5),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.True(user.HasCapability(capability.PublishArticles.String()))
	assert.False(user.HasCapability(capability.DeletePages.String()))
	assert.False(user.HasCapability(capability.CreateUsers.String()))
}
