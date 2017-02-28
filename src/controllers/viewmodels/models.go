package viewmodels

type Article struct {
	ID, Title, Slug, Content, Excerpt string
	FeaturedImage                     struct {
						  ID, Name, Caption, AltText string
						  Width, Height              int
						  Size, Url                  string
					  }
	Author                            struct {
						  ID, FullName, NickName, Website,
						  Biography, ProfilePictureUrl string
					  }
	Editor                            struct {
						  ID, FullName, NickName string
					  }
	Categories                        []struct {
		ID, Name, Slug string
	}
	Taxonomies                        []struct {
		ID, Name, Slug string
	}
	Status, CreatedAt, UpdatedAt      string
}

type Page struct {
	ID, Title, Slug, Content, Template,
	AuthorName, LastEditorName, Status, CreatedAt, UpdatedAt string
}

type Category struct {
	ID, Name, Slug, Description                      string
	Parent                                           struct {
								 ID, Name, Slug string
							 }
	NumArticles                                      uint16
	AuthorName, LastEditorName, CreatedAt, UpdatedAt string
}

type Taxonomy struct {
	ID, Name, Slug, Description                      string
	Parent                                           struct {
								 ID, Name, Slug string
							 }
	NumArticles                                      uint16
	AuthorName, LastEditorName, CreatedAt, UpdatedAt string
}

type Picture struct {
	ID, Name                                  string
	Width, Height                             int
	Size, Caption, Description, MimeType, Url string
	Uploader                                  struct {
							  ID, Name string
						  }
	UploadedAt, UpdatedAt                     string
}

type User struct {
	ID, FullName, NickName, Email, Website, Biography   string
	ProfilePicture                                      struct {
								    ID, Url string
							    }

	Role                                                struct {
								    ID, Name string
							    }
	CreatedBy                                           struct {
								    ID, FullName string
							    }
	NumPages, NumCategories, NumTaxonomies, NumArticles uint16
	CreatedAt, UpdatedAt                                string
}
