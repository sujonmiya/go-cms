package dashboard

import (
	"github.com/gorilla/mux"
	"services"
	"net/http"
	"controllers/viewmodels"
	"models"
	"middleware/auth"
	"github.com/satori/go.uuid"
	"repository"
)

type MediaLibraryController struct {
	r       *mux.Router
	s       *services.PictureService
}

func NewMediaLibraryController(r *mux.Router) *MediaLibraryController {
	return &MediaLibraryController{
		r: r.StrictSlash(true).PathPrefix("/media-library").Subrouter(),
		s: services.NewPictureService(),
	}
}

func (mlc *MediaLibraryController) RegisterEndpoints() {
	/*mlc.r.Path("/").
		Methods(http.MethodGet).
		Queries("new", "{new}").
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(ac.NewArticleHandler))
		HandlerFunc(mlc.NewArticleHandler)

	mlc.r.Path("/{id}").
		Methods(http.MethodGet).
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(ac.UpdateArticleHandler))
		HandlerFunc(mlc.UpdateArticleHandler)*/

	mlc.r.Path("/").
		Methods(http.MethodGet).
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(ac.ArticlesHandler))
		HandlerFunc(mlc.MediaLibraryHandler)
}

func (mlc *MediaLibraryController) MediaLibraryHandler(w http.ResponseWriter, r *http.Request) {
	/*p, err := auth.GetUserPrincipal(r)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}*/

	var data struct {
		Pictures    []*viewmodels.Picture
		CurrentUser *models.UserPrincipal
		CsrfToken   string
		Signature   string
		Error       string
	}

	//data.CurrentUser = user
	data.CsrfToken = auth.CsrfToken(r)
	data.Signature = uuid.NewV4().String()

	pics, err := mlc.s.GetPicturesByQuery(repository.NewDefaultQuery())
	if err != nil {
		data.Error = err.Error()
		renderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	data.Pictures = pics
	renderer.HTML(w, http.StatusOK, "media-library", &data)
}