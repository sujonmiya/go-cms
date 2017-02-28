package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/icrowley/fake"
	"golang.org/x/crypto/bcrypt"
	"utils/seed"
	"repository"
)

var (
	userService *UserService
)

func init() {
	_service := NewUserService()
	userService = _service
}

func TestUserService(t *testing.T) {
	ass := assert.New(t)
	ass.NotNil(userService)
}

func TestUserService_UserExist(t *testing.T) {
	ass := assert.New(t)
	exist := userService.UserExist(fake.EmailAddress())
	ass.False(exist)
}

func TestUserService_EncryptPassword(t *testing.T) {
	ass := assert.New(t)
	pass := []byte(fake.SimplePassword())
	encrypted, err := encryptPassword(pass)
	ass.NoError(err)
	ass.Len(encrypted, 60)
	err = bcrypt.CompareHashAndPassword(encrypted, pass)
	ass.NoError(err)
}


func TestCreateUser(t *testing.T) {
	ass := assert.New(t)
	newUser := seed.NewAdministrator()
	user, err := userService.CreateUser(newUser)
	ass.NotNil(user)
	ass.NoError(err)
	user, err = userService.CreateUser(newUser)
	ass.Error(err)
	ass.Nil(user)
}

func TestUserService_GetUsersByQuery(t *testing.T) {
	ass := assert.New(t)
	users, err := userService.GetUsersByQuery(repository.NewDefaultQuery())
	ass.NoError(err)
	ass.Len(users, 2)
}

/*
func TestGetUser(t *testing.T) {
	assert := assert.New(t)

	user, err := userService.GetUser(user.Id.Hex())
	assert.NotNil(user)
	assert.NoError(err)
}

func TestGetUserByEmail(t *testing.T) {
	assert := assert.New(t)

	user, err := userService.GetUserByEmail(user.Email)
	assert.NotNil(user)
	assert.NoError(err)
}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)

	err := userService.UpdateUser(user)
	assert.NoError(err)
}

func TestGetUsers(t *testing.T) {
	assert := assert.New(t)

	users, err := userService.GetUsers(utils.NewDefaultQuery())
	assert.NotEmpty(users)
	assert.NoError(err)
}

func TestVerifyLogin(t *testing.T) {
	assert := assert.New(t)

	user, err := userService.VerifyLogin(user.Email, password)
	assert.NotNil(user)
	assert.NoError(err)
}

func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)

	err := userService.DeleteUser(user.Id.Hex())
	assert.NoError(err)
}

func TestGetUserNotFound(t *testing.T) {
	assert := assert.New(t)

	user, err := userService.GetUser(user.Id.Hex())
	assert.Equal(mgo.ErrNotFound, err)
	assert.Nil(user)
}*/
