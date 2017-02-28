package api

import (
	"services"
	"github.com/gorilla/mux"
	"net/http"
	"models"
)

type PictureApiController struct {
	r       *mux.Router
	service *services.PictureService
}

func NewPictureApiController(r *mux.Router) *PictureApiController {
	return &PictureApiController{
		r: r,
		service: services.NewPictureService(),
	}
}

func (pac *PictureApiController) RegisterEndpoints() {
	pac.r.Methods(http.MethodGet).
		HeadersRegexp(HeaderAccept, MediaTypeJson).
		//Name(capabilities.UpdateCategory.String()).
		//Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(pac.ListHandler))
		HandlerFunc(pac.ListHandler)

	pac.r.Methods(http.MethodPost).
		HeadersRegexp(HeaderAccept, MediaTypeJson,
		HeaderContentType, MediaTypeJson).
		//Name(capabilities.UpdateCategory.String()).
		//Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(pac.CreateHandler))
		HandlerFunc(pac.CreateHandler)

	pac.r.Path("/{id}").
		Methods(http.MethodPut).
		HeadersRegexp(HeaderAccept, MediaTypeJson,
		HeaderContentType, MediaTypeJson).
		//Name(capabilities.UpdateCategory.String()).
		//Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(pac.EditHandler))
		HandlerFunc(pac.EditHandler)

	pac.r.Path("/{id}").
		Methods(http.MethodDelete).
		//Name(capabilities.UpdateCategory.String()).
		//Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(pac.DeleteHandler))
		HandlerFunc(pac.DeleteHandler)
}

func (pac *PictureApiController) ListHandler(w http.ResponseWriter, r *http.Request) {
	/*var query models.Query
	if err := ParseForm(r, &query); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}*/

	pictures, err := pac.service.GetPictures()
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusOK, pictures)
}

func (pac *PictureApiController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var picture models.NewPicture
	if err := ParseJson(r.Body, &picture); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if err := picture.Validate(); err != nil {
		renderer.JSON(w, 422, models.ValidationErr(err))
		return
	}

	/*var c models.NewCategoryRequest
	if err := ParseJson(r.Body, &c); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
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
		parent, err := pac.service.GetCategory(c.Parent.Hex())
		if err != nil {
			renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
			return
		}

		p := converters.Category2Taxonomy(parent)
		category.Parent = &p
	}*/

	saved, err := pac.service.CreateAndSavePicture(&picture)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusCreated, saved)
}

func (pac *PictureApiController) EditHandler(w http.ResponseWriter, r *http.Request) {
	/*var c models.UpdateCategoryRequest
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
		parent, err := pac.service.GetCategory(c.Parent.Hex())
		if err != nil {
			renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
			return
		}

		category.Parent = converters.Category2Taxonomy(parent)
	}

	if err := pac.service.UpdateCategory(&category); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}*/

	renderer.JSON(w, http.StatusOK, &models.Picture{})
}

func (pac *PictureApiController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	/*vars := mux.Vars(r)
	category, err := pac.service.GetCategory(vars["id"])
	if err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if err := pac.service.DeleteCategory(category.Id.Hex()); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}*/

	renderer.JSON(w, http.StatusNoContent, nil)
}
