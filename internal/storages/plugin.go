package storages

import (
	"github.com/arwos/artifactory/internal/proxy/db"
	"github.com/arwos/artifactory/internal/storages/files"
	"github.com/arwos/artifactory/internal/storages/storage"
	"github.com/arwos/artifactory/internal/storages/users"
	"github.com/deweppro/goppy/plugins"
)

var Plugins = plugins.Plugins{}.Inject(
	db.Plugin,
	users.Plugin,
	storage.Plugin,
	files.Plugin,
)
