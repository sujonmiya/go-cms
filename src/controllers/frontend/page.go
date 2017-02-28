package frontend

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"services"
	"controllers"
	"github.com/thoas/stats"
	"controllers/viewmodels"
)

type Page struct {
	Title     string
	Content   string
	UpdatedAt time.Time
}

type PageController struct {
	r       *mux.Router
	service *services.PageService
}

func NewPageController(r *mux.Router) *PageController {
	return &PageController{
		r:       r,
		service: services.NewPageService(),
	}
}

func (c *PageController) RegisterEndpoints() {
	c.r.Path("/").
		Methods(http.MethodGet).
		HandlerFunc(c.ShowHomePage)
	c.r.Path("/page/{slug}").
		Methods(http.MethodGet).
		HandlerFunc(c.ShowPage)

	c.r.Path("/stats").
		Methods(http.MethodGet).
		HandlerFunc(c.ShowStats)
}

func (c *PageController) ShowPage(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
		Page  *viewmodels.Page
		Error string
	}

	vars := mux.Vars(r)
	page, err := c.service.GetPage(vars["slug"])
	if err != nil {
		data.Title = http.StatusText(http.StatusNotFound)
		data.Error = err.Error()
		controllers.Renderer.HTML(w, http.StatusNotFound, "404", &data)
		return
	}

	data.Page = page
	controllers.Renderer.HTML(w, http.StatusOK, "page", &data)
}

func (c *PageController) ShowHomePage(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
		Articles  []Article
		Error string
	}

	/*vars := mux.Vars(r)
	page, err := c.service.GetPageBySlug(vars["slug"])
	if err != nil {
		data.Title = http.StatusText(http.StatusNotFound)
		data.Error = err.Error()
		controllers.Renderer.HTML(w, http.StatusNotFound, "404", &data)
		return
	}

	data.Title = page.Title
	data.Page = Page{
		Title:     page.Title,
		Content:   page.Content,
		UpdatedAt: page.UpdatedAt,
	}

	template := strings.ToLower(page.Template)
	if template == "" || template == "default" {
		template = "page"
	}*/
	data.Articles = []Article{}

	controllers.Renderer.HTML(w, http.StatusOK, "index", &data)
}

func (c *PageController) ShowStats(w http.ResponseWriter, r *http.Request) {

	/*vars := mux.Vars(r)
	page, err := c.service.GetPageBySlug(vars["slug"])
	if err != nil {
		data.Title = http.StatusText(http.StatusNotFound)
		data.Error = err.Error()
		controllers.Renderer.HTML(w, http.StatusNotFound, "404", &data)
		return
	}

	data.Title = page.Title
	data.Page = Page{
		Title:     page.Title,
		Content:   page.Content,
		UpdatedAt: page.UpdatedAt,
	}

	template := strings.ToLower(page.Template)
	if template == "" || template == "default" {
		template = "page"
	}*/

	data := stats.New().Data()
	controllers.Renderer.JSON(w, http.StatusOK, &data)
}
