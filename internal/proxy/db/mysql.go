package db

import (
	"github.com/deweppro/go-sdk/orm"
	"github.com/deweppro/goppy/plugins"
	"github.com/deweppro/goppy/plugins/database"
)

var Plugin = plugins.Plugin{
	Inject: newDB,
}

type (
	dbProxy struct {
		db database.MySQL
	}

	DB interface {
		Main() orm.Stmt
	}
)

func newDB(db database.MySQL) DB {
	return &dbProxy{
		db: db,
	}
}

func (v *dbProxy) Main() orm.Stmt {
	return v.db.Pool("main")
}
