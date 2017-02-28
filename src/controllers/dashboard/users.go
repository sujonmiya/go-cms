package dashboard

import (
	"net/http"

	"github.com/gorilla/mux"
	"services"
	"middleware/auth"
	"models"
	"github.com/satori/go.uuid"
	"repository"
	"controllers/viewmodels"
)

type UsersController struct {
	r       *mux.Router
	service *services.UserService
}

func NewUsersController(r *mux.Router) *UsersController {
	return &UsersController{
		r: r.StrictSlash(true).PathPrefix("/users").Subrouter(),
		service: services.NewUserService(),
	}
}

func (uc *UsersController) RegisterEndpoints() {
	uc.r.Path("/").
		Methods(http.MethodGet).
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(uc.UsersHandler))
		HandlerFunc(uc.UsersHandler)

	uc.r.Path("/{id}").
		Methods(http.MethodGet).
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(uc.NewUserHandler))
		HandlerFunc(uc.NewUserHandler)

	uc.r.Path("/{id}").
		Methods(http.MethodPut).
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(uc.UpdateUserHandler))
		HandlerFunc(uc.UpdateUserHandler)
}

func (c *UsersController) UsersHandler(w http.ResponseWriter, r *http.Request) {
	/*user, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}*/

	var data struct {
		Users       []*viewmodels.User
		Roles       []models.Role
		CurrentUser *models.UserPrincipal
		CsrfToken   string
		Signature   string
		Error       string
	}

	//data.CurrentUser = user
	//data.Roles = roles.Roles()
	data.CsrfToken = auth.CsrfToken(r)
	data.Signature = uuid.NewV4().String()

	users, err := c.service.GetUsersByQuery(repository.NewDefaultQuery())
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", err.Error())
		return
	}

	data.Users = users
	renderer.HTML(w, http.StatusOK, "users", &data)
}

func (c *UsersController) NewUserHandler(w http.ResponseWriter, r *http.Request) {
	p, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	data := struct {
		Roles       []string
		CurrentUser *models.UserPrincipal
		CsrfToken   string
	}{
		CurrentUser: p,
		//Roles:       roles.Roles(),
		CsrfToken:   auth.CsrfToken(r),
	}

	renderer.HTML(w, http.StatusOK, "user", &data)
}

func (c *UsersController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	p, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	var data struct {
		User        *models.User
		Roles       []string
		CurrentUser *models.UserPrincipal
		CsrfToken   string
		Error       string
	}

	data.CurrentUser = p
	//data.Roles = roles.Roles()
	data.CsrfToken = auth.CsrfToken(r)

	/*vars := mux.Vars(r)
	user, err := c.service.GetUser(vars["id"])
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", err.Error())
		return
	}

	data.User = user*/
	renderer.HTML(w, http.StatusOK, "user", &data)
}
