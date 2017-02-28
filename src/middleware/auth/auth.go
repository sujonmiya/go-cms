package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"models"
)

type key int

const (
	userKey key = iota

	tokenLength       = 64
	AuthCookieName    = "AUTH-TOKEN"
	csrfCookieName    = "XSRF-TOKEN"
	CsrfHeaderName    = "X-XSRF-TOKEN"
	csrfFormFieldName = "CsrfToken"
)

var (
	hashKey      = securecookie.GenerateRandomKey(tokenLength)
	blockKey     = securecookie.GenerateRandomKey(tokenLength / 2)
	secureCookie = securecookie.New(hashKey, blockKey)
	csrfSecret   = []byte("dsgfvdsgfdeygkhdsbhjsdbvhjsddshk")
)

type errResponse struct {
	Code    int
	Reason  string
	Message string
}

func newErrResponse(code int, err string) errResponse {
	return errResponse{
		Code:    code,
		Reason:  http.StatusText(code),
		Message: err,
	}
}

func UnauthorizedErr(err string) errResponse {
	return newErrResponse(http.StatusUnauthorized, err)
}

func Acl(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUserPrincipal(r)
		if err != nil {
			log.Printf("error current user not found in context: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(UnauthorizedErr(err.Error()))
			return
		}

		c := mux.CurrentRoute(r).GetName()
		if !HasCapability(user.Role.String(), c) {
			log.Printf("user doesn't have capability: %v", c)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(UnauthorizedErr(err.Error()))
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func HasCapability(role, capability string) bool {
	return true
}

func Auth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			log.Printf("auth coockie not present in request: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(UnauthorizedErr(err.Error()))
			return
		}

		var user models.UserPrincipal
		if err := secureCookie.Decode(AuthCookieName, cookie.Value, &user); err != nil {
			log.Printf("error decoding auth coockie: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(newErrResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		context.Set(r, userKey, &user)
		h.ServeHTTP(w, r)
	}

	//return alice.New(Cors).Then(http.HandlerFunc(fn))
	return alice.New(Cors).ThenFunc(fn)
}

func AuthRedirect(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		/*cookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			log.Printf("auth coockie not present in request: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var user models.UserPrincipal
		if err := secureCookie.Decode(AuthCookieName, cookie.Value, &user); err != nil {
			log.Printf("error decoding auth coockie: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		context.Set(r, userKey, &user)*/
		h.ServeHTTP(w, r)
	}

	return alice.New(Cors, Csrf).ThenFunc(fn)
}

func Cors(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8888"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-XSRF-TOKEN", "X-AUTH-TOKEN"},
		AllowCredentials: true,
		//Debug:true,
	})

	return c.Handler(h)
}

func SetAuthCookie(u *models.User, w http.ResponseWriter) error {
	user := models.UserPrincipal{
		//Id:       u.Id.Hex(),
		//FullName: u.FullName(),
		//Email:    u.Email,
		//Role:     u.Role.String(),
	}

	encoded, err := secureCookie.Encode(AuthCookieName, &user)
	if err != nil {
		return err
	}

	c := http.Cookie{
		Name:     AuthCookieName,
		Value:    encoded,
		Path:     "/dashboard",
		HttpOnly: true,
		//Secure:   true, //dev only
	}

	http.SetCookie(w, &c)
	c.Path = "/api/v1"
	http.SetCookie(w, &c)

	return nil
}

func ClearAuthCookie(w http.ResponseWriter) {
	c := http.Cookie{
		Name:   AuthCookieName,
		Value:  "",
		MaxAge: -60 * 60 * 24,
	}

	http.SetCookie(w, &c)
}

func GetUserPrincipal(r *http.Request) (*models.UserPrincipal, error) {
	/*data, ok := context.GetOk(r, userKey)
	if !ok {
		log.Printf("user not found in context: %v", ok)
		return nil, errors.New("User not found in context")
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil {
		log.Printf("error encoding user data from context: %v", err)
		return nil, err
	}

	var user models.CurrentUser
	if err := gob.NewDecoder(&buf).Decode(&user); err != nil {
		log.Printf("error decoding user data from buffer: %v", err)
		return nil, err
	}
	buf.Reset()*/

	cookie, err := r.Cookie(AuthCookieName)
	if err != nil {
		log.Printf("auth coockie not present in request: %v", err)
		return nil, err
	}
	var user models.UserPrincipal
	if err := secureCookie.Decode(AuthCookieName, cookie.Value, &user); err != nil {
		log.Printf("error decoding auth coockie: %v", err)
		return nil, err
	}

	return &user, nil
}

func Csrf(h http.Handler) http.Handler {
	fn := csrf.Protect(csrfSecret,
		csrf.FieldName(csrfFormFieldName),
		csrf.CookieName(csrfCookieName),
		csrf.RequestHeader(CsrfHeaderName),
		csrf.Path("/"),
		csrf.HttpOnly(true),
		csrf.Secure(false), //dev only
		csrf.ErrorHandler(http.HandlerFunc(csrfErrorHandler)))

	return fn(h)
}

func CsrfToken(r *http.Request) string {
	return csrf.Token(r)
}

func csrfErrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("error in Csrf: %v", csrf.FailureReason(r))
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(newErrResponse(http.StatusForbidden, csrf.FailureReason(r).Error()))
	return
}
