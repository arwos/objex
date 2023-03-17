package main

import (
	"github.com/arwos/artifactory/internal/controllers"
	"github.com/arwos/artifactory/internal/providers"
	"github.com/arwos/artifactory/internal/proxy"
	"github.com/arwos/artifactory/internal/storages"
	"github.com/deweppro/goppy"
	"github.com/deweppro/goppy/plugins/database"
	"github.com/deweppro/goppy/plugins/web"
)

func main() {
	app := goppy.New()
	app.WithConfig("./config.yaml") // Reassigned via the `--config` argument when run via the console.
	app.Plugins(
		web.WithHTTP(),
		web.WithHTTPClient(),
		database.WithMySQL(),
	)
	app.Plugins(providers.Plugin)
	app.Plugins(proxy.Plugins...)
	app.Plugins(storages.Plugins...)
	app.Plugins(controllers.Plugins...)
	app.Run()
}
