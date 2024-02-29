package db

import (
	"go.osspkg.com/goppy/orm"
	"go.osspkg.com/goppy/ormmysql"
)

type (
	object struct {
		db ormmysql.MySQL
	}

	DB interface {
		Master() orm.Stmt
		Slave() orm.Stmt
	}
)

func New(db ormmysql.MySQL) DB {
	return &object{
		db: db,
	}
}

func (v *object) Master() orm.Stmt {
	return v.db.Pool("main")
}

func (v *object) Slave() orm.Stmt {
	return v.db.Pool("main")
}
