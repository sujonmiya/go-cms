package main

import (
	"log"
	"runtime"
	"net/http"
	"github.com/gorilla/mux"
	"controllers"
	"controllers/dashboard"
	"controllers/frontend"
	"os"
	"github.com/urfave/cli"
	"fmt"
	"time"
	"github.com/urfave/negroni"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/thoas/stats"
	"github.com/meatballhat/negroni-logrus"
	"github.com/rs/cors"
	"gopkg.in/tylerb/graceful.v1"
	"services"
	"utils/seed"
	"models"
)

var (
	debugEnabled bool
	configFileName string
	seedDatabase bool
	runServer bool
	serverAddress string
	serverPort int
	verboseMode bool
)

const (
	serverCommandName = "server"
	createCommandName = "create"
)

func SeedDummy() {
	userService := services.NewUserService()
	pageService := services.NewPageService()
	categoryService := services.NewCategoryService()
	taxonomyService := services.NewTaxonomyService()
	articleService := services.NewArticleService()
	user, err := userService.CreateUser(seed.NewAdministrator())
	if err != nil {
		log.Printf("Could not seed the database: %v", err)
		return
	}

	userPrincipal := &models.UserPrincipal{ID:user.ID}
	for i := 0; i < 5; i++ {
		page := seed.NewPage()
		page.Author = userPrincipal
		_, err = pageService.SavePage(page)
		if err != nil {
			continue
		}
	}

	for i := 0; i < 5; i++ {
		category := seed.NewCategory()
		category.Author = userPrincipal
		c, err := categoryService.SaveCategory(category)
		if err != nil {
			continue
		}

		taxonomy := seed.NewTaxonomy()
		taxonomy.Author = userPrincipal
		t, err := taxonomyService.SaveTaxonomy(taxonomy)
		if err != nil {
			continue
		}

		article := seed.NewArticle()
		article.Author = userPrincipal
		article.Categories = append(article.Categories, c.ID)
		article.Taxonomies = append(article.Taxonomies, t.ID)
		for i := 0; i < 2; i++ {
			_, err = articleService.SaveArticle(article)
			if err != nil {
				continue
			}
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	/***
	CHANGE config base path !!!
	***/

	/*v1 := router.Group("/api/v1")

	pages := v1.Group("/pages")
	pages.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})
	pages.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})
	pages.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, id)
	})
	pages.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, id)
	})
	pages.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, id)
	})

	articles := v1.Group("/articles")
	articles.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})
	articles.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})
	articles.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, id)
	})
	articles.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, id)
	})
	articles.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, id)
	})

	categories := v1.Group("/categories")
	categories.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})
	categories.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})
	categories.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, id)
	})
	categories.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, id)
	})
	categories.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, id)
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"Title": "Login"})
	})
	router.GET("/sign-up", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sign-up.html", gin.H{"Title": "Login"})
	})
	router.GET("/forgot-password", func(c *gin.Context) {
		c.HTML(http.StatusOK, "forgot-password.html", gin.H{"Title": "Login"})
	})

	dashboard := router.Group("dashboard")
	dashboard.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{"Title": "Dashboard"})
	})
	dashboard.GET("/pages", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard/pages.tmpl", gin.H{"Title": "Dashboard"})
	})
	dashboard.GET("/articles", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard/articles.tmpl", gin.H{"Title": "Dashboard"})
	})
	dashboard.GET("/media-library", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard/media-library.tmpl", gin.H{"Title": "Dashboard"})
	})
	dashboard.GET("/categories", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard/categories.tmpl", gin.H{"Title": "Dashboard"})
	})
	dashboard.GET("/users", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard/users.tmpl", gin.H{"Title": "Dashboard"})
	})
	dashboard.GET("/themes", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard/themes.tmpl", gin.H{"Title": "Dashboard"})
	})
	dashboard.GET("/settings", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard/settings.tmpl", gin.H{"Title": "Dashboard"})
	})
	dashboard.GET("/help", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard/help.tmpl", gin.H{"Title": "Dashboard"})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"Title": "Home"})
	})
	router.GET("/page/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		c.HTML(http.StatusOK, "page.html", gin.H{"Title": slug})
	})

	router.GET("/article/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		c.HTML(http.StatusOK, "article.html", gin.H{"Title": slug})
	})

	router.GET("/category/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		c.HTML(http.StatusOK, "category.html", gin.H{"Title": slug})
	})*/

	router := mux.NewRouter()
	router.PathPrefix("/assets/").
		Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./templates/assets"))))
	router.PathPrefix("/uploads/").
		Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	login := dashboard.NewLoginController(router)
	login.RegisterEndpoints()
	sign := dashboard.NewSignUpController(router)
	sign.RegisterEndpoints()

	dcr := router.StrictSlash(true).PathPrefix("/dashboard").Subrouter()
	dashboard.NewDashboardController(dcr).RegisterEndpoints()

	frontend.NewPageController(router).RegisterEndpoints()
	frontend.NewArticlesController(router).RegisterEndpoints()
	frontend.NewCategoriesController(router).RegisterEndpoints()

	v1 := router.StrictSlash(true).PathPrefix("/api/v1").Subrouter()
	pages := v1.PathPrefix("/pages").Subrouter()
	controllers.NewPageApiController(pages).RegisterEndpoints()

	articles := v1.PathPrefix("/articles").Subrouter()
	controllers.NewArticleApiController(articles).RegisterEndpoints()


	app := cli.NewApp()
	app.Name = "cms"
	app.Usage = "The CMS for Contetto"
	app.Version = "1.0.0"
	app.Authors = []cli.Author{{
		Name:  "Sujon Miya",
		Email: "dhongsabshesh@gmail.com",
	},}
	app.Copyright = fmt.Sprintf("(c) %d Sujon Miya", time.Now().Year())
	cli.HelpFlag = cli.BoolFlag{
		Name: "help,h",
		Usage: "Shows a list of commands or help for one command",}
	cli.VersionFlag = cli.BoolFlag{
		Name: "version,v",
		Usage: "Prints version information",}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "debug,d",
			Usage: "Enable debug mode",
			Destination: &debugEnabled,
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
			Destination: &configFileName,
		},
		cli.BoolFlag{
			Name: "seed,s",
			Usage: "Seed the database with dummy data",
			Destination: &seedDatabase,
		},
	}

	app.Commands = []cli.Command{{
			Name: serverCommandName,
			Aliases: []string{"S"},
			Usage: "The CMS server",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "run,r",
					Usage:"Run the server",
					Destination: &runServer,
				},
				cli.StringFlag{
					Name: "addr,a",
					Usage: "Address to listen on",
					Destination: &serverAddress,
				},
				cli.IntFlag{
					Name: "port,p",
					Value: 3000,
					EnvVar: "CMS_PORT",
					Usage: "Port to connect to",
					Destination: &serverPort,
				},
				cli.BoolFlag{
					Name: "verbose,V",
					Usage:"Verbose mode",
					Destination: &verboseMode,
				},
			},
			Action:func(c *cli.Context) error {
				if c.NumFlags() == 0 && c.NArg() == 0 {
					return cli.ShowCommandHelp(c, serverCommandName)
				}

				if c.Bool("run") {
					n := negroni.New(negroni.NewRecovery())
					if verboseMode {
						n.Use(negronilogrus.NewMiddleware())
					}

					n.Use(cors.Default())
					n.Use(gzip.Gzip(gzip.DefaultCompression))
					n.Use(stats.New())

					/*vg := vangoh.New()
					_ = vg.AddProvider("API", services.NewUserService())
					n.Use(negroni.HandlerFunc(vg.NegroniHandler))*/

					n.UseHandler(router)
					addrPortFormat := fmt.Sprintf("%s:%d", serverAddress, serverPort)
					log.Printf("[CMS-debug] Server listening and serving HTTP on %s\n", addrPortFormat)
					return graceful.RunWithErr(addrPortFormat, 10 * time.Second, n)
				}

				return nil
			},
		},
		{
			Name: createCommandName,
			Aliases: []string{"C"},
			Usage: "Create entities on the database",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "user,u",
					Usage:"Create a new User",
				},
			},
			Action:func(c *cli.Context) error {
				if c.NumFlags() == 0 && c.NArg() == 0 {
					return cli.ShowCommandHelp(c, createCommandName)
				}

				if c.Bool("user") {
					user, err := services.NewUserService().CreateUser(seed.NewAdministrator())
					if err != nil {
						return err
					}

					log.Printf("%+v", user)
				}

				return nil
			},
		},}

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 && c.NumFlags() == 0 {
			return cli.ShowAppHelp(c)
		}

		switch true {
		case debugEnabled:
			log.Println("log debug messages")
		case seedDatabase:
			log.Println("Seed the database please")
			SeedDummy()
		}

		return nil
	}

	app.Run(os.Args)
}
