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

type ArticleApiController struct {
	*mux.Router
	service *services.ArticleService
}

func NewArticleApiController(r *mux.Router) *ArticleApiController {
	return &ArticleApiController{
		Router:r,
		service:services.NewArticleService(),
	}
}

func (aac *ArticleApiController) RegisterEndpoints() {
	aac.Methods(http.MethodGet).
		HeadersRegexp(HeaderAccept, MediaTypeJson).
		Name(capabilities.ReadArticle.String()).
		HandlerFunc(aac.listHandler)

	aac.Methods(http.MethodPost).
		HeadersRegexp(HeaderAccept, MediaTypeJson,
			HeaderContentType, MediaTypeJson).
		Name(capabilities.CreateArticle.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(aac.createHandler))

	aac.Path("/{id}").
		Methods(http.MethodGet).
		HeadersRegexp(HeaderAccept, MediaTypeJson).
		Name(capabilities.ReadArticle.String()).
		HandlerFunc(aac.Get)

	aac.Path("/{id}").
		Methods(http.MethodPut).
		HeadersRegexp(HeaderAccept, MediaTypeJson,
			HeaderContentType, MediaTypeJson).
		Name(capabilities.UpdateArticle.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(aac.editHandler))

	aac.Path("/{id}").
		Methods(http.MethodDelete).
		Name(capabilities.DeleteArticle.String()).
		Handler(alice.New(auth.Auth, auth.Acl).ThenFunc(aac.deleteHandler))

}

func (api *ArticleApiController) listHandler(w http.ResponseWriter, r *http.Request) {
	var query models.Query
	if err := ParseForm(r, &query); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	articles, err := api.service.GetArticles(query)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusOK, articles)
}

func (api *ArticleApiController) createHandler(w http.ResponseWriter, r *http.Request) {
	var a models.NewArticleRequest
	if err := ParseJson(r.Body, &a); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if _, err := a.Validate(); err != nil {
		renderer.JSON(w, 422, models.ValidationErr(err))
		return
	}

	user, err := auth.GetUserPrincipal(r)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	article := models.NewArticle{
		Title:      a.Title,
		Content:    a.Content,
		//Metadata:   a.Metadata,
		Status:     a.Status,
		//Visibility: a.Visibility,
		ScheduleAt: a.ScheduleAt,
		Author:     user,
		//Editor:     user,
	}

	/*categories, err := api.service.GetCategories(a.Categories)
	if err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	for _, c := range categories {
		article.Categories = append(article.Categories, converters.Category2Taxonomy(c))
	}

	if len(a.Tags) > 0 {
		tags, err := api.service.GetTags(a.Tags)
		if err != nil {
			Renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
			return
		}

		for _, t := range tags {
			article.Tags = append(article.Tags, converters.Tag2Taxonomy(t))
		}
	}*/

	saved, err := api.service.SaveArticle(&article)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	renderer.JSON(w, http.StatusCreated, saved)
}

func (api *ArticleApiController) editHandler(w http.ResponseWriter, r *http.Request) {
	var a models.UpdateArticleRequest
	if err := ParseJson(r.Body, &a); err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if _, err := a.Validate(); err != nil {
		renderer.JSON(w, 422, models.ValidationErr(err))
		return
	}

	user, err := auth.GetUserPrincipal(r)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	article := models.UpdateArticle{
		Title:      a.Title,
		Content:    a.Content,
		//Metadata:   a.Metadata,
		Status:     a.Status,
		//Visibility: a.Visibility,
		ScheduleAt: a.ScheduleAt,
		Editor:     user,
	}
	/*vars := mux.Vars(r)
	article.SetId(vars["id"])*/

	/*categories, err := api.service.GetCategories(a.Categories)
	if err != nil {
		renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	for _, c := range categories {
		article.Categories = append(article.Categories, converters.Category2Taxonomy(c))
	}

	if len(a.Tags) > 0 {
		tags, err := api.service.GetTags(a.Tags)
		if err != nil {
			Renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
			return
		}

		for _, t := range tags {
			article.Tags = append(article.Tags, converters.Tag2Taxonomy(t))
		}
	}

	if err := api.service.UpdateArticle(&article); err != nil {
		Renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}*/

	renderer.JSON(w, http.StatusOK, article)
}

func (api *ArticleApiController) Get(w http.ResponseWriter, r *http.Request) {
	/*params := mux.Vars(r)
	article, err := api.service.GetArticle(params["id"])
	if err != nil {
		Renderer.JSON(w, http.StatusNotFound, NotFoundErr(err.Error()))
		return
	}

	Renderer.JSON(w, http.StatusOK, article)*/
}

func (api *ArticleApiController) deleteHandler(w http.ResponseWriter, r *http.Request) {
	/*vars := mux.Vars(r)
	article, err := api.service.GetArticle(vars["id"])
	if err != nil {
		Renderer.JSON(w, http.StatusBadRequest, BadRequestErr(err.Error()))
		return
	}

	if err := api.service.DeleteArticle(article.Id.Hex()); err != nil {
		Renderer.JSON(w, http.StatusInternalServerError, ServerErr(err.Error()))
		return
	}

	Renderer.JSON(w, http.StatusNoContent, article)*/
}
