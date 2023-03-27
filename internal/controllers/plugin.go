package controllers

import (
	"github.com/arwos/artifactory/internal/controllers/api"
	"github.com/arwos/artifactory/internal/controllers/files"
	"github.com/arwos/artifactory/internal/controllers/npm"
	"github.com/deweppro/goppy/plugins"
)

var Plugins = plugins.Plugins{}.Inject(
	api.Plugin,
	files.Plugin,
	npm.Plugin,
)
