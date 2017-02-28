package dashboard

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"services"
	"middleware/auth"
	"models"
	"controllers"
	"code.google.com/p/go-uuid/uuid"
)

type SignUpController struct {
	r       *mux.Router
	service *services.UserService
}

func NewSignUpController(r *mux.Router) *SignUpController {
	return &SignUpController{
		r:       r,
		service: services.NewUserService(),
	}
}

func (sup *SignUpController) RegisterEndpoints() {
	sup.r.Path("/sign-up").
		Methods(http.MethodGet).
		Handler(alice.New(auth.Csrf).ThenFunc(sup.showSignUpHandler))

	sup.r.Path("/sign-up").
		Methods(http.MethodPost).
		HeadersRegexp(controllers.HeaderContentType, controllers.MediaTypeForm).
		Handler(alice.New(auth.Csrf).ThenFunc(sup.signUpHandler))
}

func (c *SignUpController) showSignUpHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
		CsrfToken string
	}{
		Title: "Sign Up",
		//CsrfToken: auth.CsrfToken(r),
		CsrfToken: uuid.New(),
	}

	renderer.HTML(w, http.StatusOK, "sign-up", &data)
}

func (c *SignUpController) signUpHandler(w http.ResponseWriter, r *http.Request) {
	var login models.Login
	if err := controllers.ParseForm(r, &login); err != nil {
		log.Printf("error parsing login form: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	/*user, err := c.service.VerifyLogin(login.Email, login.Password)
	if err != nil {
		log.Printf("error verifying login: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := auth.SetAuthCookie(user, w); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}*/

	http.Redirect(w, r, "/dashboard/", http.StatusSeeOther)
}
