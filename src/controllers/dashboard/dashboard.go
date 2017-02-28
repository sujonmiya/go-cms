package dashboard

import (
	"net/http"

	"github.com/gorilla/mux"
	"middleware/auth"
	"models"
	"github.com/unrolled/render"
	"config"
	"path/filepath"
	"html/template"
	"strings"
	"github.com/gorilla/schema"
)

var (
	renderer *render.Render
	decoder *schema.Decoder
)

func join(s []string) string {
	return strings.Join(s, ", ")
}

func init() {
	renderer = render.New(render.Options{
		Directory:     filepath.Join(config.TemplatesDir(), "dashboard"),
		Extensions:    []string{".html"},
		Funcs:         []template.FuncMap{{"join": join}},
		IsDevelopment: config.IsDevelopment(),
	})

	decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
}

type DashboardController struct {
	r *mux.Router
}

func NewDashboardController(r *mux.Router) *DashboardController {
	return &DashboardController{r:r}
}

func (dc *DashboardController) RegisterEndpoints() {
	dc.r.Path("/").
		Methods(http.MethodGet).
		//Handler(alice.New(auth.AuthRedirect).ThenFunc(dc.dashboardHandler))
		HandlerFunc(dc.dashboardHandler)

	NewPagesController(dc.r).RegisterEndpoints()
	NewArticlesController(dc.r).RegisterEndpoints()
	NewCategoriesController(dc.r).RegisterEndpoints()
	NewMediaLibraryController(dc.r).RegisterEndpoints()
	NewUsersController(dc.r).RegisterEndpoints()
}

func (c *DashboardController) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
		CurrentUser *models.UserPrincipal
		CsrfToken   string
	}{
		Title: "Dashboard",
		CsrfToken: auth.CsrfToken(r),
	}

	/*user, err := auth.GetUserPrincipal(r)
	if err != nil {
		log.Printf("sadasd: %v", err)
	} else {
		data.CurrentUser = user
	}*/

	renderer.HTML(w, http.StatusOK, "dashboard", &data)
}


