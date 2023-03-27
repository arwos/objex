package users

import (
	"context"
	"time"

	"github.com/arwos/artifactory/internal/proxy/db"
	"github.com/deweppro/go-sdk/app"
	"github.com/deweppro/go-sdk/routine"
	"github.com/deweppro/goppy/plugins"
)

var Plugin = plugins.Plugin{
	Inject: NewUsers,
}

type Users struct {
	db    db.DB
	cache *Cache
}

func NewUsers(db db.DB) *Users {
	return &Users{
		db:    db,
		cache: NewCache(1000),
	}
}

func (v *Users) Up(ctx app.Context) error {
	routine.Interval(ctx.Context(), 60*time.Minute, func(ctx context.Context) {
		v.cache.FlushAll()
	})
	return nil
}

func (v *Users) Down() error {
	return nil
}
