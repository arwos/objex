package api

import (
	"github.com/arwos/artifactory/internal/pkg/middlewares"
	"github.com/arwos/artifactory/internal/storages/users"
	"github.com/deweppro/goppy/plugins"
	"github.com/deweppro/goppy/plugins/web"
)

var Plugin = plugins.Plugin{
	Config: &Config{},
	Inject: NewController,
}

type Controller struct {
	conf   *Config
	users  *users.Users
	routes web.RouterPool
}

func NewController(r web.RouterPool, u *users.Users, c *Config) *Controller {
	return &Controller{
		users:  u,
		routes: r,
		conf:   c,
	}
}

func (v *Controller) Up() error {
	route := v.routes.Main()

	v.InjectUIRoutes(route)

	v.InjectAuthRoutes(route.Collection("/api/auth", middlewares.TokenDetectMiddleware(v.conf.Settings.CookieName)))

	//route.Post("/files/{storage}/#", v.UploadFile)
	//route.Get("/files/{storage}/#", v.DownloadFile)
	//
	//apiV1 := route.Collection("/files/api/v1")
	//apiV1.Post("users/new", v.CreateUser)
	//apiV1.Post("users/group/add", v.AddUserGroup)
	//apiV1.Post("groups/new", v.CreateGroup)
	//apiV1.Get("groups/list", v.ListGroup)
	//apiV1.Post("storage/new", v.CreateStorage)
	//apiV1.Post("storage/group/add", v.AddStorageGroup)
	//apiV1.Post("search/props", v.SearchByProps)

	return nil
}

func (v *Controller) Down() error {
	return nil
}
