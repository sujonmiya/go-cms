package repository

import (
	"models"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	"config"
	"log"
	"strings"
	"fmt"
	"github.com/markbates/pop/nulls"
	"models/roles"
	"time"
)

const (
	dialect = "mysql"
	cascade = "CASCADE"
	restrict = "RESTRICT"
	TableViewArticles = "vw_articles"
)

var (
	db *gorm.DB
)

type Repository struct{}

type Filter struct {
	Author   []uint32
	Editor   []uint32
	Category []uint32
	Taxonomy []uint32
	Status   string
	Total    uint16
	Offset   uint8
	Sort     string `schema:"-"`
}

/*func (f Filter) OrderBy() string {
	return strings.Join(f.Sort, " ")
}*/

func NewRepo() *Repository {
	return &Repository{}
}

type ArticleCategory struct {
	ArticleID  nulls.UInt32 `sql:"index"`
	CategoryID nulls.UInt32 `sql:"index"`
}

type ArticleTaxonomy struct {
	ArticleID  nulls.UInt32 `sql:"index"`
	TaxonomyID nulls.UInt32 `sql:"index"`
}

type UserRole struct {
	UserID nulls.UInt32 `sql:"index"`
	RoleID nulls.UInt32 `sql:"index"`
}

type UserLogin struct {
	UserID                uint32 `sql:"index"`
	IPAddress, ClientName string
	Time                  time.Time
}

type ArticleResult struct {
	ArticleID                                                                    uint32
	Title, Slug, Content, Excerpt                                                string
	FeaturedImageID                                                              uint32
	FeaturedImageName, FeaturedImageCaption, FeaturedImageDesc, FeaturedImageUrl string
	FeaturedImageWidth, FeaturedImageHeight, FeaturedImageSize                   int
	AuthorID                                                                     uint32
	AuthorName, AuthorNickName, AuthorWebsite, AuthorBio, AuthorProfilePicUrl    string
	EditorID                                                                     uint32
	EditorName, EditorNickName,
	Categories                                                                   string
	Taxonomies                                                                   string
	LastEditorName, Status                                                       string
	CreatedAt, UpdatedAt                                                         time.Time
}

type PageResult struct {
	PageID                                 uint32
	Title, Slug, Content, Template, Status string
	AuthorID                               uint32
	AuthorName, LastEditorName             string
	CreatedAt, UpdatedAt                   time.Time
}

type CategoryResult struct {
	CategoryID                      uint32
	Name, Slug, Description, Parent string
	NumArticles                     uint16
	AuthorID                        uint32
	AuthorName                      string
	LastEditorName                  string
	CreatedAt, UpdatedAt            time.Time
}

type PictureResult struct {
	PictureID                       uint32
	Name                            string
	Width, Height, Size             int
	Caption, AltText, MimeType, Url string
	UploaderID                      uint32
	UploaderName                    string
	CreatedAt, UpdatedAt            time.Time
}

type UserResult struct {
	UserID                                              uint32
	FullName, NickName, Email, Website, Biography       string
	ProfilePicID                                        uint32
	ProfilePicUrl                                       string
	RoleID                                              uint32
	RoleName                                            string
	CreatedByID                                         uint32
	CreatedByName                                       string
	NumPages, NumCategories, NumTaxonomies, NumArticles uint16
	CreatedAt, UpdatedAt                                time.Time
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dsn := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		config.DbUsername(), config.DbPassword(), config.Dbname())
	_db, err := gorm.Open(dialect, dsn)
	if err != nil {
		panic(err.Error())
	}

	db = _db
	db.LogMode(config.IsDevelopment())

	if !db.HasTable(&models.Role{}) {
		db.CreateTable(&models.Role{})
		tx := db.Begin()
		for _, r := range roles.Roles() {
			role := &models.Role{Name:r.String()}
			if err := tx.Create(role).Error; err != nil {
				tx.Rollback()
				log.Fatalf("Error creating Role: %v", err)
			}
		}

		tx.Commit()
	}

	// DB Migration order matters. DO NOT Alter!!
	db.AutoMigrate(&models.User{}).
		AddForeignKey("created_by_id", "users(id)", restrict, cascade).
		AddForeignKey("role_id", "roles(id)", restrict, cascade)
	db.AutoMigrate(&UserRole{}).
		AddForeignKey("user_id", "users(id)", restrict, cascade).
		AddForeignKey("role_id", "roles(id)", restrict, cascade)
	db.AutoMigrate(&UserLogin{}).AddForeignKey("user_id", "users(id)", restrict, cascade)

	db.AutoMigrate(&models.Picture{}).
		AddForeignKey("author_id", "users(id)", restrict, cascade).
		AddForeignKey("last_editor_id", "users(id)", restrict, cascade)
	db.Model(&models.User{}).AddForeignKey("profile_picture_id", "pictures(id)", restrict, cascade)

	db.AutoMigrate(&models.Page{}).
		AddForeignKey("author_id", "users(id)", restrict, cascade).
		AddForeignKey("last_editor_id", "users(id)", restrict, cascade)

	db.AutoMigrate(&models.Article{}).
		AddForeignKey("picture_id", "pictures(id)", restrict, cascade).
		AddForeignKey("author_id", "users(id)", restrict, cascade).
		AddForeignKey("last_editor_id", "users(id)", restrict, cascade)

	db.AutoMigrate(&models.Category{}).
		AddForeignKey("parent_id", "categories(id)", restrict, cascade).
		AddForeignKey("author_id", "users(id)", restrict, cascade).
		AddForeignKey("last_editor_id", "users(id)", restrict, cascade)

	db.AutoMigrate(&ArticleCategory{}).
		AddForeignKey("article_id", "articles(id)", restrict, cascade).
		AddForeignKey("category_id", "categories(id)", restrict, cascade)

	db.AutoMigrate(&models.Taxonomy{}).
		AddForeignKey("parent_id", "taxonomies(id)", restrict, cascade).
		AddForeignKey("author_id", "users(id)", restrict, cascade).
		AddForeignKey("last_editor_id", "users(id)", restrict, cascade)

	db.AutoMigrate(&ArticleTaxonomy{}).
		AddForeignKey("article_id", "articles(id)", restrict, cascade).
		AddForeignKey("taxonomy_id", "taxonomies(id)", restrict, cascade)

	err = db.Exec(`CREATE OR REPLACE VIEW vw_articles
		  AS
		  SELECT
		  a.id                                        AS article_id,
		  a.title,
		  a.slug,
		  a.content,
		  a.excerpt,
		  a.status,
		  a.created_at,
		  a.updated_at,
		  p.id                                         AS featured_image_id,
		  p.alt_text                                   AS featured_image_desc,
		  p.caption                                    AS featured_image_caption,
		  p.url                                        AS featured_image_url,
		  u.id                                         AS author_id,
		  CONCAT_WS(' ', u.first_name, u.last_name)    AS author_name,
		  u.nick_name                                  AS author_nick_name,
		  u.website                                    AS author_website,
		  u.bio                                        AS author_bio,
		  p2.url                                       AS author_profile_pic_url,
		  u2.id                                        AS editor_id,
		  CONCAT_WS(' ', u2.first_name, u2.last_name)  AS editor_name,
		  u2.nick_name                                 AS editor_nick_name,
		  (SELECT GROUP_CONCAT(DISTINCT CONCAT_WS(':', c.id, c.name, c.slug) ORDER BY c.name)
		   FROM categories AS c
		   WHERE c.id IN (SELECT ac.category_id
				  FROM article_categories AS ac
				  WHERE ac.article_id = a.id)) AS categories,
		  (SELECT GROUP_CONCAT(DISTINCT CONCAT_WS(':', t.id, t.name, t.slug) ORDER BY t.name)
		   FROM taxonomies AS t
		   WHERE t.id IN (SELECT at.taxonomy_id
				  FROM article_taxonomies AS at
				  WHERE at.article_id = a.id)) AS taxonomies
		FROM articles AS a
		  INNER JOIN users AS u ON a.author_id = u.id
		  LEFT JOIN users AS u2 ON a.last_editor_id = u2.id
		  LEFT JOIN pictures AS p ON a.picture_id = p.id
		  LEFT JOIN pictures AS p2 ON p2.id = u.profile_picture_id
		  WHERE a.deleted_at IS NULL`).Error
	if err != nil {
		panic(err.Error())
	}
}

func (r *Repository) DB() *gorm.DB {
	return db
}

func (r *Repository) Save(doc interface{}) error {
	ch := make(chan error, 1)
	go func() {
		ch <- db.Create(doc).Error
	}()

	return <-ch
}

func (r *Repository) Update(filter interface{}, doc interface{}) error {
	ch := make(chan error, 1)
	go func() {
		ch <- db.Model(filter).
			Updates(doc).Error
	}()

	return <-ch
}

func (r *Repository) Find(dst interface{}) error {
	ch := make(chan error, 1)
	go func() {
		ch <- db.Limit(50).Find(dst).Error
	}()

	return <-ch
}

func (r *Repository) FindByQuery(q models.Query, dst interface{}) error {
	ch := make(chan error, 1)
	go func() {
		q = sanitizeQuery(q)
		ch <- db.Limit(q.Total).
			Offset(q.Offset).
			Order(q.Sort).
			Find(dst).Error
	}()

	return <-ch
}

func (r *Repository) FindByQueryAndFilter(q models.Query, filter interface{}, dst interface{}) error {
	ch := make(chan error, 1)
	go func() {
		q = sanitizeQuery(q)
		ch <- db.Limit(q.Total).
			Offset(q.Offset).
			Order(q.Sort).
			Where(filter).
			Find(dst).Error
	}()

	return <-ch
}

func (r *Repository) FindByID(id uint32, dst interface{}) error {
	ch := make(chan error, 1)
	go func() {
		ch <- db.First(dst, id).Error
	}()

	return <-ch
}

func (r *Repository) FindOne(filter interface{}, dst interface{}) error {
	ch := make(chan error, 1)
	go func() {
		ch <- db.Where(filter).
			First(dst).Error
	}()

	return <-ch
}

func (r *Repository) Delete(source interface{}) error {
	ch := make(chan error, 1)
	go func() {
		ch <- db.Delete(source).Error
	}()

	return <-ch
}

func formatSort(s string) string {
	return strings.Join(
		strings.Split(strings.Replace(s, " ", "", -1), ","), " ")
}

func NewDefaultQuery() models.Query {
	return models.NewQuery(50, 0, "created_at")
}

func sanitizeQuery(q models.Query) models.Query {
	query := models.Query{
		Total:  q.Total,
		Offset: q.Offset,
		Sort:   formatSort(q.Sort),
	}

	if query.Total < 0 || query.Total > 50 {
		query.Total = 50
	}

	if query.Offset < 0 {
		query.Offset = 0
	}

	if query.Sort == "" {
		query.Sort = "created_at"
	}

	return query
}