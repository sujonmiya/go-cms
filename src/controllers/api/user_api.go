package api

import (
	"contetto"
	"html"
	"log"
	"net/http"
	"service/app/middleware/auth"
	"service/app/models"
	"service/app/models/capability"
	"service/app/models/role"
	"service/app/services"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gopkg.in/mgo.v2/bson"
)

const (
	UserApiEndpoint = "/users"
)

type UserApi struct {
	service *services.UserService
	r       *mux.Router
}

func InitUserApi(f contetto.ContettoMicroServiceBaseFramework) *UserApi {
	log.Printf("initializing user api...")

	api := &UserApi{
		service: services.NewUserService(f),
		r: ApiBaseRouter(f.Web.Router).StrictSlash(true).
			PathPrefix(UserApiEndpoint).
			Subrouter(),
	}

	api.RegisterEndpoints()

	return api
}

func (ua *UserApi) GetRouter() *mux.Router {
	return ua.r
}

func (ua *UserApi) RegisterEndpoints() {
	log.Printf("registering user api endpoints on: %s%s", ApiBaseEndpoint, UserApiEndpoint)

	ua.r.Methods("GET").
		Name(capability.ReadUsers.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(ua.ListHandler))

	ua.r.Methods("POST").
		HeadersRegexp(HeaderAccept, MediaTypeJson,
			HeaderContentType, MediaTypeJson).
		Name(capability.CreateUsers.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(ua.CreateHandler))

	ua.r.Path("/{id}").
		Methods("PUT").
		HeadersRegexp(HeaderAccept, MediaTypeJson,
			HeaderContentType, MediaTypeJson).
		Name(capability.EditUsers.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(ua.EditHandler))

	ua.r.Path("/{id}").
		Methods("DELETE").
		Name(capability.DeleteUsers.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(ua.DeleteHandler))

	ua.r.Path("/roles").
		Methods("GET").
		Name(capability.ReadUsers.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(ua.RolesHandler))
}

func (ua *UserApi) ListHandler(w http.ResponseWriter, r *http.Request) {
	var query models.Query
	if err := ParseForm(r, &query); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	users, err := ua.service.GetUsers(query)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusOK, users)
}

func (ua *UserApi) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var user models.NewUserRequest
	if err := ParseJson(r.Body, &user); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if ok, err := user.Validate(); !ok {
		renderer.JSON(w, http.StatusBadRequest, models.ValidationErr(err))
		return
	}

	p, err := auth.GetUserPrincipal(r)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	u := models.NewUser{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		CreatedBy: p,
	}

	saved, err := ua.service.CreateUser(u)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusCreated, saved)
}

func (ua *UserApi) EditHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UpdateUserRequest
	if err := ParseJson(r.Body, &user); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if ok, err := user.Validate(); !ok {
		renderer.JSON(w, http.StatusBadRequest, models.ValidationErr(err))
		return
	}

	vars := mux.Vars(r)
	u := models.UpdateUser{
		Id:        bson.ObjectIdHex(vars["id"]),
		FirstName: user.FirstName,
		LastName:  html.EscapeString(user.LastName),
		Password:  user.Password,
		Role:      user.Role,
		Website:   user.Website,
		Bio:       user.Bio,
	}

	if err := ua.service.UpdateUser(&u); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusOK, u)
}

func (ua *UserApi) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user, err := ua.service.GetUser(vars["id"])
	if err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if err := ua.service.DeleteUser(user.Id.Hex()); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusNoContent, user)
}

func (ua *UserApi) RolesHandler(w http.ResponseWriter, r *http.Request) {
	renderer.JSON(w, http.StatusOK, role.Roles())
}
