package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"services"
	"middleware/auth"
	"models"
	"models/capabilities"
)

type PageApiController struct {
	*mux.Router
	service *services.PageService
}

func NewPageApiController(r  *mux.Router) *PageApiController {
	return &PageApiController{
		Router: r,
		service:services.NewPageService(),
	}
}

func (pac *PageApiController) RegisterEndpoints() {
	pac.Methods(http.MethodGet).
		HeadersRegexp(HeaderAccept, MediaTypeJson).
		Name(capabilities.ReadPage.String()).
		HandlerFunc(pac.ListHandler)

	pac.Methods(http.MethodPost).
		HeadersRegexp(HeaderAccept, MediaTypeJson,
			HeaderContentType, MediaTypeJson).
		Name(capabilities.CreatePage.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(pac.CreateHandler))

	pac.Path("/{id}").
		Methods(http.MethodGet).
		HeadersRegexp(HeaderAccept, MediaTypeJson).
		Name(capabilities.ReadPage.String()).
		HandlerFunc(pac.GetHandler)

	pac.Path("/{id}").
		Methods(http.MethodPut).
		HeadersRegexp(HeaderAccept, MediaTypeJson,
			HeaderContentType, MediaTypeJson).
		Name(capabilities.UpdatePage.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(pac.EditHandler))

	pac.Path("/{id}").
		Methods(http.MethodDelete).
		Name(capabilities.DeletePage.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(pac.DeleteHandler))
}

func (api *PageApiController) ListHandler(w http.ResponseWriter, r *http.Request) {
	var query models.Query
	if err := ParseForm(r, &query); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	pages, err := api.service.GetPages(query)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusOK, pages)
}

func (api *PageApiController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var p models.NewPageRequest
	if err := ParseJson(r.Body, &p); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if _, err := p.Validate(); err != nil {
		renderer.JSON(w, 422, models.ValidationErr(err))
		return
	}

	user, err := auth.GetUserPrincipal(r)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	article := models.NewPage{
		Title:      p.Title,
		Content:    p.Content,
		Template:   p.Template,
		Status:     p.Status,
		ScheduleAt: p.ScheduleAt,
		Author:     user,
	}

	saved, err := api.service.SavePage(&article)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusCreated, saved)
}

func (api *PageApiController) EditHandler(w http.ResponseWriter, r *http.Request) {
	var p models.UpdatePageRequest
	if err := ParseJson(r.Body, &p); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if _, err := p.Validate(); err != nil {
		renderer.JSON(w, 422, models.ValidationErr(err))
		return
	}

	user, err := auth.GetUserPrincipal(r)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	page := models.UpdatePage{
		Title:      p.Title,
		Content:    p.Content,
		Template:   p.Template,
		Status:     p.Status,
		Visibility: p.Visibility,
		ScheduleAt: p.ScheduleAt,
		Editor:     user,
	}
	/*vars := mux.Vars(r)
	page.SetId(vars["id"])

	if err := api.service.UpdatePage(&page); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}*/

	renderer.JSON(w, http.StatusOK, page)
}

func (api *PageApiController) GetHandler(w http.ResponseWriter, r *http.Request) {
	/*vars := mux.Vars(r)
	article, err := api.service.GetPage(vars["id"])
	if err != nil {
		renderer.JSON(w, http.StatusNotFound, NotFoundErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusOK, article)
	*/
}

func (api *PageApiController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	/*vars := mux.Vars(r)
	page, err := api.service.GetPage(vars["id"])
	if err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if err := api.service.DeletePage(page.Id.Hex()); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusNoContent, page)*/
}
