package proxy

import (
	"github.com/arwos/artifactory/internal/proxy/db"
	"github.com/deweppro/goppy/plugins"
)

var Plugins = plugins.Plugins{}.Inject(
	db.Plugin,
)
