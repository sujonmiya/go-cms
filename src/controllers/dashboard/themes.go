package dashboard

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"middleware/auth"
)

type ThemesController struct {
	r *mux.Router
}

func (c *ThemesController) RegisterEndpoints() {
	r := c.r.PathPrefix("/themes").
		Subrouter()

	r.Path("/").
		Methods("GET").
		Handler(alice.New(auth.Csrf, auth.AuthRedirect).ThenFunc(c.ThemesHandler))

}

func (c *ThemesController) ThemesHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		CsrfToken string
	}{
		CsrfToken: auth.CsrfToken(r),
	}

	renderer.HTML(w, http.StatusOK, "themes", &data)
}
