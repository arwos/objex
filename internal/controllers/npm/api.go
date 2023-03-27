package npm

import (
	"os"

	"github.com/arwos/artifactory/internal/pkg/locker"
	"github.com/arwos/artifactory/internal/pkg/network"
	"github.com/deweppro/goppy/plugins"
	"github.com/deweppro/goppy/plugins/web"
)

var Plugin = plugins.Plugin{
	Config: &Config{},
	Inject: NewController,
}

type Controller struct {
	routes web.RouterPool
	cli    network.Request
	conf   *Config
	mux    locker.Locker
}

func NewController(r web.RouterPool, conf *Config) *Controller {
	return &Controller{
		routes: r,
		cli:    network.NewRequest(),
		conf:   conf,
		mux:    locker.New(),
	}
}

func (v *Controller) Up() error {
	route := v.routes.Main()

	route.Get(registry, v.IndexYarn)
	route.Get(registry+"/#", v.LoadMetaData)
	route.Put(registry+"/#", v.PublishPackage)
	route.Get(registryFiles+"/#", v.DownloadPackage)

	users := route.Collection(registry + "/-/user")
	users.Put("/org.couchdb.user:{uid}", v.UserLogin)
	users.Delete("/token/{token}", v.DeleteToken)

	return os.MkdirAll(v.conf.Packages.ProxyCache, 0755)
}

func (v *Controller) Down() error {
	return nil
}
