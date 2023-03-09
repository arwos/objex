package main

import (
	"github.com/arwos/artifactory/internal/providers"
	"github.com/deweppro/goppy"
	"github.com/deweppro/goppy/plugins/database"
	"github.com/deweppro/goppy/plugins/web"
)

func main() {
	app := goppy.New()
	app.WithConfig("./config.yaml") // Reassigned via the `--config` argument when run via the console.
	app.Plugins(
		web.WithHTTP(),
		database.WithMySQL(),
	)
	app.Plugins(providers.Plugin)
	app.Run()
}
