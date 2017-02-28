package dashboard

import (
	"net/http"

	"github.com/gorilla/mux"
	"services"
	"middleware/auth"
	"models"
	"models/status"
	"models/visibility"
	"repository"
	"controllers/viewmodels"
	"github.com/satori/go.uuid"
	"strconv"
)

type PagesController struct {
	r       *mux.Router
	service *services.PageService
}

func NewPagesController(r *mux.Router) *PagesController {
	return &PagesController{
		r:       r.StrictSlash(true).PathPrefix("/pages").Subrouter(),
		service: services.NewPageService(),
	}
}

func (pc *PagesController) RegisterEndpoints() {
	pc.r.Methods(http.MethodGet).
		Queries("new", "{new}").
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(pc.NewPageHandler))
		HandlerFunc(pc.NewPageHandler)

	pc.r.Path("/{id:[1-9]([0-9]?)+}").
		Methods(http.MethodGet).
		Queries("task", "edit").
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(pc.UpdatePageHandler))
		HandlerFunc(pc.UpdatePageHandler)

	pc.r.Path("/{id:[1-9]([0-9]?)+}").
		Methods(http.MethodGet).
		Queries("task", "delete").
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(pc.UpdatePageHandler))
		HandlerFunc(pc.DeletePageHandler)

	pc.r.Path("/").
		Methods(http.MethodGet).
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(pc.PagesHandler))
		HandlerFunc(pc.PagesHandler)
}

func (c *PagesController) PagesHandler(w http.ResponseWriter, r *http.Request) {
	/*user, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}*/

	var data struct {
		Pages       []*viewmodels.Page
		CurrentUser *models.UserPrincipal
		CsrfToken   string
		Signature   string
		Error       string
	}

	//data.CurrentUser = user
	data.CsrfToken = auth.CsrfToken(r)
	data.Signature = uuid.NewV4().String()

	pages, err := c.service.GetPagesByQuery(repository.NewDefaultQuery())
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", data)
		return
	}

	data.Pages = pages
	renderer.HTML(w, http.StatusOK, "pages", &data)
}

func (c *PagesController) NewPageHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	var data struct {
		CurrentUser  *models.UserPrincipal
		Templates    []string
		Statuses     []string
		Visibilities []string
		CsrfToken    string
	}

	data.CurrentUser = user
	data.Statuses = status.Statuses()
	data.Visibilities = visibility.Visibilities()
	data.CsrfToken = auth.CsrfToken(r)

	renderer.HTML(w, http.StatusOK, "page-new", &data)
}

func (c *PagesController) DeletePageHandler(w http.ResponseWriter, r *http.Request) {

}

func (c *PagesController) UpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	/*user, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}*/

	var data struct {
		Page         *viewmodels.Page
		CurrentUser  *models.UserPrincipal
		Templates    []string
		Statuses     []string
		Visibilities []string
		CsrfToken    string
		Error        string
	}

	//data.CurrentUser = user
	data.Statuses = status.Statuses()
	data.Visibilities = visibility.Visibilities()
	data.CsrfToken = auth.CsrfToken(r)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", err.Error())
		return
	}

	page, err := c.service.GetPageByID(uint32(id))
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", err.Error())
		return
	}

	data.Page = page
	renderer.HTML(w, http.StatusOK, "page-edit", &data)
}
