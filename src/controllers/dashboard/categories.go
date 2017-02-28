package dashboard

import (
	"net/http"

	"github.com/gorilla/mux"
	"services"
	"middleware/auth"
	"models"
	"controllers/viewmodels"
	"repository"
	"github.com/satori/go.uuid"
	"strconv"
)

type CategoriesController struct {
	r *mux.Router
	s *services.CategoryService
}

func NewCategoriesController(r *mux.Router) *CategoriesController {
	return &CategoriesController{
		r: r.StrictSlash(true).PathPrefix("/categories").Subrouter(),
		s: services.NewCategoryService(),
	}
}

func (cc *CategoriesController) RegisterEndpoints() {
	cc.r.Path("/").
		Methods(http.MethodGet).
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(cc.CategoriesHandler))
		HandlerFunc(cc.CategoriesHandler)

	cc.r.Path("/{id:[1-9]([0-9]?)+}").
		Methods(http.MethodGet).
		Queries("task", "edit").
	//Handler(alice.New(auth.AuthRedirect).ThenFunc(cc.UpdateCategoryHandler))
		HandlerFunc(cc.UpdateCategoryHandler)

	cc.r.Path("/{id:[1-9]([0-9]?)+}").
		Methods(http.MethodGet).
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(cc.NewCategoryHandler))
		HandlerFunc(cc.NewCategoryHandler)
}

func (cc *CategoriesController) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	/*user, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}*/

	var data struct {
		Categories  []*viewmodels.Category
		CurrentUser *models.UserPrincipal
		CsrfToken   string
		Signature   string
		Error       string
	}

	//data.CurrentUser = user
	data.CsrfToken = auth.CsrfToken(r)
	data.Signature = uuid.NewV4().String()
	categories, err := cc.s.GetCategoriesByQuery(repository.NewDefaultQuery())
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", err.Error())
		return
	}

	data.Categories = categories
	renderer.HTML(w, http.StatusOK, "categories", &data)
}

func (c *CategoriesController) NewCategoryHandler(w http.ResponseWriter, r *http.Request) {
	/*user, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}*/

	var data struct {
		CurrentUser *models.UserPrincipal
		Parents     []*models.Category
		CsrfToken   string
		Error       string
	}

	//data.CurrentUser = user
	data.CsrfToken = auth.CsrfToken(r)

	/*parents, err := c.service.GetCategoriesFromQuery(utils.NewDefaultQuery())
	if err != nil {
		data.Error = err.Error()
		DashboardRenderer.HTML(w, http.StatusOK, "error", err.Error())
		return
	}

	data.Parents = parents*/
	renderer.HTML(w, http.StatusOK, "category", &data)
}

func (c *CategoriesController) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Category    *viewmodels.Category
		CurrentUser *models.UserPrincipal
		CsrfToken   string
		Error       string
	}

	/*user, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	data.CurrentUser = user*/

	data.CsrfToken = auth.CsrfToken(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", err.Error())
		return
	}

	category, err := c.s.GetCategoryByID(uint32(id))
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", err.Error())
		return
	}

	data.Category = category
	renderer.HTML(w, http.StatusOK, "category-edit", &data)
}
