package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/justinas/alice"
	"gopkg.in/mgo.v2/bson"
	"services"
	"middleware/auth"
	"models"
	"models/capability"
)

const (
	CategoryApiEndpoint = "/categories"
)

type CategoryApiController struct {
	r       *mux.Router
	service *services.CategoryService
}

func NewCategoryApiController(r *mux.Router) *CategoryApiController {
	return &CategoryApiController{
		r: r,
		service: services.NewCategoryService(),
	}
}

func (cac *CategoryApiController) RegisterEndpoints() {
	cac.r.Methods(http.MethodGet).
		HeadersRegexp(HeaderAccept, MediaTypeJson).
		Name(capabilities.UpdateCategory.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(cac.ListHandler))

	cac.r.Methods(http.MethodPost).
		HeadersRegexp(HeaderAccept, MediaTypeJson,
		HeaderContentType, MediaTypeJson).
		Name(capabilities.UpdateCategory.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(cac.CreateHandler))

	cac.r.Path("/{id}").
		Methods(http.MethodPut).
		HeadersRegexp(HeaderAccept, MediaTypeJson,
		HeaderContentType, MediaTypeJson).
		Name(capabilities.UpdateCategory.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(cac.EditHandler))

	cac.r.Path("/{id}").
		Methods(http.MethodDelete).
		Name(capabilities.UpdateCategory.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(cac.DeleteHandler))
}

func (api *CategoryApiController) ListHandler(w http.ResponseWriter, r *http.Request) {
	var query models.Query
	if err := ParseForm(r, &query); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	articles, err := api.service.GetCategoriesByQuery(query)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusOK, articles)
}

func (api *CategoryApiController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var c models.NewCategoryRequest
	if err := ParseJson(r.Body, &c); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if ok, err := c.Validate(); !ok {
		renderer.JSON(w, 422, models.ValidationErr(err))
		return
	}

	user, err := auth.GetUserPrincipal(r)
	if err != nil {
		renderer.JSON(w, http.StatusUnauthorized, auth.UnauthorizedErr(err.Error()))
		return
	}

	category := models.NewCategory{
		Name:        c.Name,
		Description: c.Description,
		Author:      user,
	}

	if c.Parent.Valid() {
		parent, err := api.service.GetCategory(c.Parent.Hex())
		if err != nil {
			renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
			return
		}

		p := converters.Category2Taxonomy(parent)
		category.Parent = &p
	}

	saved, err := api.service.SaveCategory(&category)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusCreated, saved)
}

func (api *CategoryApiController) EditHandler(w http.ResponseWriter, r *http.Request) {
	var c models.UpdateCategoryRequest
	if err := ParseJson(r.Body, &c); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if ok, err := c.Validate(); !ok {
		renderer.JSON(w, 422, models.ValidationErr(err))
		return
	}

	user, err := auth.GetUserPrincipal(r)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	params := mux.Vars(r)
	category := models.UpdateCategory{
		Id:          bson.ObjectIdHex(params["id"]),
		Name:        c.Name,
		Slug:        slug.Make(c.Name),
		Description: c.Description,
		Editor:      user,
	}

	if c.Parent.Valid() {
		parent, err := api.service.GetCategory(c.Parent.Hex())
		if err != nil {
			renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
			return
		}

		category.Parent = converters.Category2Taxonomy(parent)
	}

	if err := api.service.UpdateCategory(&category); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusOK, category)
}

func (api *CategoryApiController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category, err := api.service.GetCategory(vars["id"])
	if err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if err := api.service.DeleteCategory(category.Id.Hex()); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusNoContent, category)
}
