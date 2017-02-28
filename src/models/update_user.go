package models

import (
	"encoding/json"
	"net/mail"
	"net/url"
	"time"

	"gopkg.in/mgo.v2/bson"
	"models/roles"
)

type UpdateUser struct {
	Id             bson.ObjectId `bson:"-"`
	FirstName      string        `bson:"firstname"`
	LastName       string        `bson:"lastname,omitempty"`
	Password       []byte        `bson:"password,omitempty"`
	Role           roles.Role     `bson:"role"`
	Website        string        `bson:"website,omitempty"`
	Bio            string        `bson:"bio,omitempty"`
	ProfilePicture *NewPicture   `bson:",omitempty"`
	UpdatedAt      time.Time
}

type UpdateUserRequest struct {
	FirstName string
	LastName  string `json:",omitempty"`
	Password  []byte
	Role      roles.Role
	Website   string `json:",omitempty"`
	Bio       string `json:",omitempty"`
}

func (u UpdateUserRequest) Validate() (bool, *Errs) {
	var errors Errs
	p := string(u.Password)

	switch {
	case IsEmpty(u.FirstName):
		errors.AddEmptyError("FirstName")

	case LengthNotBetween(u.FirstName, 2, 64):
		errors.AddLengthBetweenError(2, 64, "FirstName", u.FirstName)

	case ContainsNotAllowed(u.FirstName, PatternAlphaSpace):
		errors.AddNotAllowedError("FirstName", u.FirstName, PatternAlphaSpace)

	case LengthExceedsMax(u.LastName, 64):
		errors.AddMaxLengthError(64, "LastName", u.LastName)

	case !IsEmpty(p) && LengthNotBetween(p, 6, 16):
		errors.AddLengthBetweenError(6, 16, "Password", u.FirstName)

	case roles.IsNotValid(u.Role):
		errors.AddInvalidError("Role", u.Role.String())
	}

	if len(errors) > 0 {
		return false, &errors
	}

	return true, nil
}

type Website string

func (w Website) MarshalJSON() ([]byte, error) {
	url, err := url.Parse(w.String())
	if err != nil {
		return nil, err
	}

	return json.Marshal(url.String())
}

func (w Website) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	w = Website(s)

	return nil
}

func (w Website) String() string {
	return string(w)
}

type Email string

func (e Email) MarshalJSON() ([]byte, error) {
	address, err := mail.ParseAddress(e.String())
	if err != nil {
		return nil, err
	}

	return json.Marshal(address.String())
}

func (e Email) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	e = Email(s)

	return nil
}

func (e Email) String() string {
	return string(e)
}
