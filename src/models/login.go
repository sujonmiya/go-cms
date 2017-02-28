package models

type Login struct {
	Email     string `schema:"Email" form:"email" json:"email" binding:"required"`
	Password  []byte `schema:"Password" form:"password" json:"password" binding:"required"`
	CsrfToken string `json:"-" schema:"CsrfToken"`
}
