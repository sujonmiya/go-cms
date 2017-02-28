package models

import (
	"fmt"
	"strings"
	"github.com/markbates/pop/nulls"
)

type User struct {
	Model
	FirstName        string `sql:"not null"`
	LastName         string
	NickName         string
	Email            string        `json:",omitempty" sql:"not null;index;unique"`
	Password         []byte        `json:"-" sql:"not null;size:60"`
	Website          string        `json:",omitempty"`
	Biography        string        `json:",omitempty" sql:"type:TEXT"`
	ProfilePicture   *Picture
	ProfilePictureID nulls.UInt32
	RoleID           nulls.UInt32 `sql:"not null;index"`
	CreatedByID      nulls.UInt32
}

func (u User) FullName() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", u.FirstName, u.LastName))
}

/*func (u User) HasCapability(s string) bool {
	c := capability.FromRouteName(s)
	if allowed, ok := capability.Acl[c]; ok {
		for _, r := range allowed {
			if r == u.Role {
				return true
			}
		}
	}

	return false
}

func (u User) capabilities() []capability.Capability {
	var caps []capability.Capability
	for c, allowed := range capability.Acl {
		for _, r := range allowed {
			if r == u.Role {
				caps = append(caps, c)
			}
		}
	}

	return caps
}

func (u User) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	buf.WriteByte(byte(u.ID))

	buf.WriteString(u.FirstName)
	buf.WriteString(u.LastName)
	buf.WriteString(u.Email)

	r, _ := u.Role.MarshalJSON()
	buf.Write(r)

	return buf.Bytes(), nil
}

func (u *User) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := fmt.Fscan(buf, &u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role)

	return err
}*/
