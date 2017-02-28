package dashboard

import (
	"net/http"

	"github.com/gorilla/mux"
	"services"
	"middleware/auth"
	"models"
	"models/status"
	"models/visibility"
	"controllers/viewmodels"
	"repository"
	"github.com/satori/go.uuid"
	"strconv"
	"log"
)

type ArticlesController struct {
	r       *mux.Router
	service *services.ArticleService
	cs      *services.CategoryService
}

func NewArticlesController(r *mux.Router) *ArticlesController {
	return &ArticlesController{
		r:       r.StrictSlash(true).PathPrefix("/articles").Subrouter(),
		service: services.NewArticleService(),
		cs:services.NewCategoryService(),
	}
}

func (ac *ArticlesController) RegisterEndpoints() {
	ac.r.Path("/").
		Methods(http.MethodGet).
		Queries("new", "{new}").
	//Handler(alice.New(auth.AuthRedirect).ThenFunc(ac.NewArticleHandler))
		HandlerFunc(ac.NewArticleHandler)

	ac.r.Path("/").
		Methods(http.MethodGet).
	//Handler(alice.New(auth.AuthRedirect).ThenFunc(ac.ArticlesHandler))
		HandlerFunc(ac.ArticlesHandler)

	ac.r.Path("/{id:[1-9]([0-9]?)+}").
		Methods(http.MethodPost).
		Queries("task", "edit").
	//Handler(alice.New(auth.AuthRedirect).ThenFunc(ac.UpdateArticleHandler))
		HandlerFunc(ac.UpdateArticleHandler)

	ac.r.Path("/{id:[1-9]([0-9]?)+}").
		Methods(http.MethodGet).
		Queries("task", "edit").
	//Handler(alice.New(auth.AuthRedirect).ThenFunc(ac.UpdateArticleHandler))
		HandlerFunc(ac.ShowUpdateArticleHandler)

	ac.r.Path("/{id:[1-9]([0-9]?)+}").
		Methods(http.MethodGet).
		Queries("task", "delete").
	//Handler(alice.New(auth.AuthRedirect).ThenFunc(ac.UpdateArticleHandler))
		HandlerFunc(ac.DeleteArticleHandler)
}

func (c *ArticlesController) ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	/*user, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}*/

	var data struct {
		Articles    []*viewmodels.Article
		CurrentUser *models.UserPrincipal
		CsrfToken   string
		Signature   string
		Error       string
	}

	var filter repository.Filter
	if err := decoder.Decode(&filter, r.URL.Query()); err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	log.Printf("%+v", filter)

	//data.CurrentUser = user
	data.CsrfToken = auth.CsrfToken(r)
	data.Signature = uuid.NewV4().String()

	articles, err := c.service.GetArticlesByFilter(filter)
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	data.Articles = articles
	renderer.HTML(w, http.StatusOK, "articles", &data)
}

func (c *ArticlesController) NewArticleHandler(w http.ResponseWriter, r *http.Request) {
	/*user, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}*/

	var data struct {
		Categories   []*models.Category
		Tags         []*models.Taxonomy
		Statuses     []string
		Visibilities []string
		CurrentUser  *models.UserPrincipal
		CsrfToken    string
		Error        string
	}

	//data.CurrentUser = user
	data.Statuses = status.Statuses()
	data.Visibilities = visibility.Visibilities()
	data.CsrfToken = auth.CsrfToken(r)

	/*categories, err := c.service.GetAllCategories()
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	data.Categories = categories
	tags, err := c.service.GetAllTags()
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	data.Tags = tags*/
	renderer.HTML(w, http.StatusOK, "article-new", &data)
}

func (c *ArticlesController) ShowUpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	/*user, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}*/

	var data struct {
		CurrentUser  *models.UserPrincipal
		Article      *viewmodels.Article
		Categories   []*viewmodels.Category
		Tags         []*viewmodels.Taxonomy
		Statuses     []string
		Visibilities []string
		CsrfToken    string
		Signature    string
		Error        string
	}

	//data.CurrentUser = user
	data.Statuses = status.Statuses()
	data.Visibilities = visibility.Visibilities()
	data.CsrfToken = auth.CsrfToken(r)
	data.Signature = uuid.NewV4().String()

	/*params := mux.Vars(r)
	article, err := c.service.GetArticle(params["id"])
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "404", &data)
		return
	}

	data.Article = article
	categories, err := c.service.GetAllCategories()
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	data.Categories = categories
	tags, err := c.service.GetAllTags()
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	data.Tags = tags*/
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	article, err := c.service.GetArticleByID(uint32(id))
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	data.Article = article

	cats, err := c.cs.GetCategories()
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}
	data.Categories = cats

	renderer.HTML(w, http.StatusOK, "article-edit", &data)
}

func (c *ArticlesController) UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	/*user, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}*/

	var data struct {
		CurrentUser  *models.UserPrincipal
		Article      *viewmodels.Article
		Statuses     []string
		Visibilities []string
		CsrfToken    string
		Signature    string
		Error        string
	}

	//data.CurrentUser = user
	data.Statuses = status.Statuses()
	data.Visibilities = visibility.Visibilities()
	data.CsrfToken = auth.CsrfToken(r)
	data.Signature = uuid.NewV4().String()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	article, err := c.service.GetArticleByID(uint32(id))
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	if err := r.ParseForm(); err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	var art models.UpdateArticleRequest
	if err := decoder.Decode(&art, r.PostForm); err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	log.Printf("%+v", art)


	data.Article = article

	renderer.HTML(w, http.StatusOK, "article-edit", &data)
}

func (c *ArticlesController) DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {

}
