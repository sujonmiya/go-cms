package capabilities

//go:generate stringer -type=Capability
//go:generate jsonenums -type=Capability
type Capability uint8

const (
	_ = iota
	CreateUser Capability = iota
	CreatePage
	CreateArticle
	CreateCategory
	CreateTaxonomy
	UploadPicture
	InstallTheme

	ReadUser
	ReadPage
	ReadArticle
	ReadCategory
	ReadTaxonomy
	ReadPicture
	ReadTheme

	UpdateUser
	UpdatePage
	UpdateArticle
	UpdateCategory
	UpdateTaxonomy
	UpdatePicture
	SwitchTheme

	DeleteUser
	DeletePage
	DeleteArticle
	DeleteCategory
	DeleteTaxonomy
	DeletePicture
	DeleteTheme

	ManageConfigs
)

func Capabilities() []Capability {
	var capas []Capability
	for _, c := range _CapabilityNameToValue {
		capas = append(capas, c)
	}

	return capas
}

func FromRouteName(s string) Capability {
	c, ok := _CapabilityNameToValue[s]
	if !ok {
		return ReadArticle
	}

	return c
}
