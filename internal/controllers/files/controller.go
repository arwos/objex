package files

import (
	"github.com/arwos/artifactory/internal/providers"
	"github.com/arwos/artifactory/internal/storages/files"
	"github.com/arwos/artifactory/internal/storages/storage"
	"github.com/arwos/artifactory/internal/storages/users"
	"github.com/deweppro/goppy/plugins"
	"github.com/deweppro/goppy/plugins/web"
)

var Plugin = plugins.Plugin{
	Inject: NewController,
}

type Controller struct {
	users     *users.Users
	providers providers.Providers
	routes    web.RouterPool
	store     *storage.Storages
	files     *files.Files
}

func NewController(
	r web.RouterPool, u *users.Users, p providers.Providers,
	s *storage.Storages, f *files.Files,
) *Controller {
	return &Controller{
		users:     u,
		providers: p,
		routes:    r,
		store:     s,
		files:     f,
	}
}

func (v *Controller) Up() error {
	route := v.routes.Main()

	route.Post("/files/{storage}/#", v.UploadFile)
	route.Get("/files/{storage}/#", v.DownloadFile)

	apiV1 := route.Collection("/files/api/v1")
	apiV1.Post("users/new", v.CreateUser)
	apiV1.Post("users/group/add", v.AddUserGroup)
	apiV1.Post("groups/new", v.CreateGroup)
	apiV1.Get("groups/list", v.ListGroup)
	apiV1.Post("storage/new", v.CreateStorage)
	apiV1.Post("storage/group/add", v.AddStorageGroup)
	apiV1.Post("search/props", v.SearchByProps)

	return nil
}

func (v *Controller) Down() error {
	return nil
}
