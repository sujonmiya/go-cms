package config

import (
	"models/events"
	"path/filepath"
	"os"
	"io/ioutil"
	"models"
	"fmt"
	"strconv"
	"log"
	"gopkg.in/yaml.v2"
	"github.com/vrischmann/envconfig"
	"time"
	"github.com/kardianos/osext"
	"models/roles"
	"models/capabilities"
)

var (
	basePath string
	conf *config
	Acl map[roles.Role][]capabilities.Capability
)

const (
	envPrefix string = "CMS"
	themesDirName = "themes"
	templatesDirName = "templates"
	uploadsDirName = "uploads"
	themeConfigFile = "theme.yaml"
)

type config struct {
	IsProduction    bool   `envconfig:"default=false"`
	ActiveTheme     string `envconfig:"default=default"`
	Port            string `envconfig:"default=:8080"`
	Database        struct {
		                Username string `envconfig:"default=root"`
		                Password string `envconfig:"default=sujon"`
		                Name     string `envconfig:"default=cms"`
	                }
	MailServer      struct {
		                Host     string
		                Port     int
		                Login    string
		                Password string
	                } `envconfig:"-"`
	Cookie          struct {
		                HashKey  string
		                BlockKey string
	                } `envconfig:"-"`
	Events          chan events.Event `envconfig:"-"`
}

func init() {
	var _conf config
	err := envconfig.InitWithPrefix(&_conf, envPrefix)
	if err != nil {
		panic(err)
	}

	conf = &_conf
	_baseDir, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}

	basePath = _baseDir
	basePath = `C:\Users\Sujon\Documents\Projects\Contetto\cms`
	log.Printf("basePath: %s\n", basePath)

	dir := filepath.Join(basePath, uploadsDirName)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0664); err != nil {
				panic(err)
			}
		}

		panic(err)
	}

	conf.Events = make(chan events.Event)

	go func() {
		for {
			event := <-conf.Events
			switch event {
			case events.ThemeChanged:
				log.Printf("ThemeChanged: %v", event)

			case events.DatabaseNameChanged:
				log.Printf("DatabaseNameChanged: %v", event)

			case events.MailServerHostChanged:
				log.Printf("MailServerHostChanged: %v", event)

			case events.MailServerLoginChanged:
				log.Printf("MailServerLoginChanged: %v", event)

			case events.MailServerPortChanged:
				log.Printf("MailServerPortChanged: %v", event)

			case events.MailServerPasswordChanged:
				log.Printf("MailServerPasswordChanged: %v", event)
			}
		}
	}()

	Acl = make(map[roles.Role][]capabilities.Capability)
	Acl[roles.Subscriber] = []capabilities.Capability{capabilities.ReadPage, capabilities.ReadArticle}

	//Acl[roles.Editor] = append(Acl[roles.Editor], Acl[roles.Subscriber]...)
	Acl[roles.Editor] = append(Acl[roles.Editor],
		capabilities.UpdatePage,
		capabilities.UpdateArticle,
		capabilities.ReadCategory,
		capabilities.UpdateCategory,
		capabilities.ReadTaxonomy,
		capabilities.UpdateTaxonomy,
		capabilities.ReadPicture,
		capabilities.UpdatePicture)

	//Acl[roles.Author] = append(Acl[roles.Author], Acl[roles.Editor]...)
	Acl[roles.Author] = append(Acl[roles.Author],
		capabilities.CreatePage,
		capabilities.DeletePage,
		capabilities.CreateArticle,
		capabilities.DeleteArticle,
		capabilities.CreateCategory,
		capabilities.DeleteCategory,
		capabilities.CreateTaxonomy,
		capabilities.DeleteTaxonomy,
		capabilities.UploadPicture,
		capabilities.DeletePicture,)

	//Acl[roles.Administrator] = append(Acl[roles.Administrator], Acl[roles.Author]...)
	Acl[roles.Administrator] = append(Acl[roles.Administrator],
		capabilities.CreateUser,
		capabilities.ReadUser,
		capabilities.UpdateUser,
		capabilities.DeleteUser,
		capabilities.InstallTheme,
		capabilities.ReadTheme,
		capabilities.SwitchTheme,
		capabilities.DeleteTheme,
		capabilities.ManageConfigs,)
}

func Port() string {
	return conf.Port
}

func IsDevelopment() bool {
	return !conf.IsProduction
}

func Config() *config {
	return conf
}

func DbUsername() string {
	return conf.Database.Username
}

func DbPassword() string {
	return conf.Database.Password
}

func Dbname() string {
	return conf.Database.Name
}

func EnvKey(field string) string {
	return fmt.Sprintf("%s%s", envPrefix, field)
}

func ThemesDir() string {
	return filepath.Join(basePath, themesDirName)
}

func TemplatesDir() string {
	return filepath.Join(basePath, templatesDirName)
}

func ActiveThemeTemplatesDir() string {
	return filepath.Join(basePath, themesDirName, conf.ActiveTheme, templatesDirName)
}

func BasePath() string {
	return basePath
}

func UploadsDir() (absolutePath string, relativePath string) {
	now := time.Now().UTC()
	relativePath = fmt.Sprintf("/%s/%d/%d/%d", uploadsDirName, now.Year(), now.Month(),
		now.Day())
	absolutePath = filepath.Join(basePath, relativePath)
	return
}

func (c *config) SetActiveTheme(name string) error {
	c.ActiveTheme = name
	c.Events <- events.ThemeChanged
	return os.Setenv(EnvKey("ACTIVE_THEME"), name)
}

func (c *config) SetDatabaseName(name string) error {
	c.Database.Name = name
	c.Events <- events.DatabaseNameChanged
	return os.Setenv(EnvKey("DATABASE_NAME"), name)
}

func (c *config) SetMailServerHost(h string) error {
	c.MailServer.Host = h
	c.Events <- events.MailServerHostChanged
	return os.Setenv(EnvKey("MAIL_SERVER_HOST"), h)
}

func (c *config) SetMailServerPort(p int) error {
	c.MailServer.Port = p
	c.Events <- events.MailServerPortChanged
	return os.Setenv(EnvKey("MAIL_SERVER_PORT"), strconv.Itoa(p))
}

func (c *config) SetMailServerLogin(l string) error {
	c.MailServer.Login = l
	c.Events <- events.MailServerLoginChanged
	return os.Setenv(EnvKey("MAIL_SERVER_LOGIN"), l)
}

func (c *config) SetMailServerPassword(p string) error {
	c.MailServer.Password = p
	c.Events <- events.MailServerPasswordChanged
	return os.Setenv(EnvKey("MAIL_SERVER_PASSWORD"), p)
}

func ThemesAssetsDir() string {
	return fmt.Sprintf("%s/%s/assets", ThemesDir(), conf.ActiveTheme)
}

func GetInstalledThemes() []models.Theme {
	var themes []models.Theme
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}

	path := filepath.Join(dir, "themes")
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("error reading themes dir: %v", err)
		return themes
	}

	for _, info := range fileInfo {
		if info.IsDir() {
			path := filepath.Join(path, info.Name(), themeConfigFile)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("error reading theme config file: %v", err)
				continue
			}

			var theme models.Theme
			if err := yaml.Unmarshal(data, &theme); err != nil {
				log.Printf("eror unmarshaling yaml file: %v", err)
				continue
			}

			themes = append(themes, theme)
		}
	}

	return themes
}
