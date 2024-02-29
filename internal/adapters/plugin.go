package adapters

import (
	"go.arwos.org/objex/internal/adapters/db"
	"go.arwos.org/objex/internal/adapters/filestore"
	"go.osspkg.com/goppy/plugins"
)

var Plugins = plugins.Inject(
	db.New,
	filestore.Plugin,
)
