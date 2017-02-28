package frontend

import (
	"net/http"

	"github.com/gorilla/mux"
	"services"
	"controllers"
	"controllers/viewmodels"
)

type CategoriesController struct {
	r *mux.Router
	*services.CategoryService
	*services.ArticleService
}

func NewCategoriesController(r *mux.Router) *CategoriesController {
	return &CategoriesController{
		r: r.StrictSlash(true).PathPrefix("/categories").Subrouter(),
		CategoryService: services.NewCategoryService(),
		ArticleService:  services.NewArticleService(),
	}
}

func (cc *CategoriesController) RegisterEndpoints() {
	cc.r.Path("/{slug}").
		Methods(http.MethodGet).
		HandlerFunc(cc.ShowCategory)
	cc.r.Path("/{category-slug}/article/{article-slug}").
		Methods(http.MethodGet).
		HandlerFunc(cc.ShowCategoryArticle)
}

func (cc *CategoriesController) ShowCategory(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title    string
		Category struct {
			         Name        string
			         Description string
			         Articles    []*viewmodels.Article
		         }
		Error    string
	}{}

	/*vars := mux.Vars(r)
	category, err := c.GetCategoryForSlug(vars["slug"])
	if err != nil {
		data.Title = http.StatusText(http.StatusNotFound)
		data.Error = err.Error()
		controllers.Renderer.HTML(w, http.StatusNotFound, "404", &data)
		return
	}

	data.Title = fmt.Sprintf("Category - %s", category.Name)
	data.Category.Name = category.Name
	data.Category.Slug = category.Slug
	data.Category.Description = category.Description

	articles, err := c.GetArticlesFromCategorySlug(category.Slug)
	if err != nil {
		data.Title = http.StatusText(http.StatusInternalServerError)
		data.Error = err.Error()
		ThemeRenderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	for _, article := range articles {
		a := struct {
			Article
			Author Author
		}{
			Article: Article{
				Title:     article.Title,
				Slug:      article.Slug,
				Content:   article.Content,
				UpdatedAt: article.UpdatedAt,
			},
			Author: Author{
				FullName: article.Author.FullName,
				Website:  article.Author.Website,
				Bio:      article.Author.Bio,
			},
		}

		data.Category.Articles = append(data.Category.Articles, a)
	}*/

	controllers.Renderer.HTML(w, http.StatusOK, "category", &data)
}

func (cc *CategoriesController) ShowCategoryArticle(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title    string
		Category struct {
			         Category struct{}
			         Articles []struct {
				         Article
				         Author Author
			         }
		         }
		Error    string
	}{}

	/*vars := mux.Vars(r)
	category, err := c.GetCategoryForSlug(vars["slug"])
	if err != nil {
		data.Title = http.StatusText(http.StatusNotFound)
		data.Error = err.Error()
		controllers.Renderer.HTML(w, http.StatusNotFound, "404", &data)
		return
	}

	data.Title = fmt.Sprintf("Category - %s", category.Name)
	data.Category.Name = category.Name
	data.Category.Slug = category.Slug
	data.Category.Description = category.Description

	articles, err := c.GetArticlesFromCategorySlug(category.Slug)
	if err != nil {
		data.Title = http.StatusText(http.StatusInternalServerError)
		data.Error = err.Error()
		ThemeRenderer.HTML(w, http.StatusOK, "error", &data)
		return
	}

	for _, article := range articles {
		a := struct {
			Article
			Author Author
		}{
			Article: Article{
				Title:     article.Title,
				Slug:      article.Slug,
				Content:   article.Content,
				UpdatedAt: article.UpdatedAt,
			},
			Author: Author{
				FullName: article.Author.FullName,
				Website:  article.Author.Website,
				Bio:      article.Author.Bio,
			},
		}

		data.Category.Articles = append(data.Category.Articles, a)
	}*/

	controllers.Renderer.HTML(w, http.StatusOK, "category", &data)
}
