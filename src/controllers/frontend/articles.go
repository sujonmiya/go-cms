package frontend

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"services"
	"models"
	"controllers"
	"repository"
	"controllers/viewmodels"
)

type ArticlesController struct {
	r       *mux.Router
	service *services.ArticleService
}

type Author struct {
	FullName string
	Website  string
	Bio      string
}

type Article struct {
	Title     string
	Slug      string
	Content   string
	Metadata  models.Metadata
	UpdatedAt time.Time
}

func NewArticlesController(r *mux.Router) *ArticlesController {
	return &ArticlesController{
		r:       r.StrictSlash(true).PathPrefix("/articles").Subrouter(),
		service: services.NewArticleService(),
	}
}

func (ac *ArticlesController) RegisterEndpoints() {
	ac.r.Path("/").
		Methods(http.MethodGet).
		HandlerFunc(ac.ShowArticles)

	ac.r.Path("/{slug}").
		Methods(http.MethodGet).
		HandlerFunc(ac.ShowArticle)
}

func (c *ArticlesController) ShowArticles(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title    string
		Articles []*viewmodels.Article
		Error    string
	}{
		Title: "Articles",
	}

	articles, err := c.service.GetArticlesByQuery(repository.NewDefaultQuery())
	if err != nil {
		data.Title = http.StatusText(http.StatusInternalServerError)
		data.Error = err.Error()
		controllers.Renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	data.Articles = articles
	controllers.Renderer.HTML(w, http.StatusOK, "articles", &data)
}

func (c *ArticlesController) ShowArticle(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title   string
		Article *viewmodels.Article
		Error   string
	}

	vars := mux.Vars(r)
	article, err := c.service.GetArticleBySlug(vars["slug"])
	if err != nil {
		data.Title = http.StatusText(http.StatusNotFound)
		data.Error = err.Error()
		controllers.Renderer.HTML(w, http.StatusNotFound, "404", &data)
		return
	}

	data.Title = article.Title
	data.Article = article
	controllers.Renderer.HTML(w, http.StatusOK, "article", &data)
}
