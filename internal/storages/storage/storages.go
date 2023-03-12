package storage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/arwos/artifactory/internal/proxy/db"
	"github.com/deweppro/go-sdk/app"
	"github.com/deweppro/go-sdk/orm"
	"github.com/deweppro/go-sdk/routine"
	"github.com/deweppro/goppy/plugins"
)

var Plugin = plugins.Plugin{
	Inject: NewStorages,
}

type Storages struct {
	db    db.DB
	cache *Cache
}

func NewStorages(db db.DB) *Storages {
	return &Storages{
		db:    db,
		cache: NewCache(1000),
	}
}

func (v *Storages) Up(ctx app.Context) error {
	routine.Interval(ctx.Context(), 60*time.Minute, func(ctx context.Context) {
		v.cache.FlushAll()
	})
	return nil
}

func (v *Storages) Down() error {
	return nil
}

func (v *Storages) Get(ctx context.Context, name string) (*Store, error) {
	if s, ok := v.cache.Get(name); ok {
		return &s, nil
	}
	s, err := v.reloadStoreFromDB(ctx, name)
	return s, err
}

func (v *Storages) reloadStoreFromDB(ctx context.Context, name string) (*Store, error) {
	var (
		sid, lifetime int64
		code          string
	)
	err := v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `id`, `lifetime`, `code` FROM `storage` WHERE `name` = ? LIMIT 1;", name)
		q.Bind(func(bind orm.Scanner) error {
			return bind.Scan(&sid, &lifetime, &code)
		})
	})
	if err != nil {
		return nil, err
	}
	if sid == 0 {
		return nil, fmt.Errorf("storage not found")
	}

	s := Store{
		ID:       sid,
		Lifetime: lifetime,
		Name:     name,
		Code:     code,
		Groups:   make(map[int64]struct{}, 0),
	}

	err = v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `group_id` FROM `storage_group` WHERE `storage_id` = ? LIMIT 1;", s.ID)
		q.Bind(func(bind orm.Scanner) error {
			var gid int64
			if err = bind.Scan(&gid); err != nil {
				return err
			}
			s.Groups[gid] = struct{}{}
			return nil
		})
	})

	v.cache.Set(s)
	return &s, nil
}

func (v *Storages) CreateStore(ctx context.Context, name, code string, lifetime int64) error {
	return v.db.Main().ExecContext("", ctx, func(q orm.Executor) {
		q.SQL("INSERT INTO `storage` (`name`, `lifetime`, `code`, `created_at`, `updated_at`) "+
			"VALUES (?, ?, ?, now(), now());", name, lifetime, code)
		q.Bind(func(result orm.Result) error {
			if result.RowsAffected == 0 {
				return fmt.Errorf("failed to create store")
			}
			return nil
		})
	})
}

func (v *Storages) DeleteStoreFromGroup(ctx context.Context, name string, groups ...int64) error {
	s, err := v.Get(ctx, name)
	if err != nil {
		return err
	}

	return v.db.Main().ExecContext("", ctx, func(q orm.Executor) {
		q.SQL("DELETE FROM `storage_group` WHERE `storage_id` = ? AND `group_id` = ?;")

		for _, gid := range groups {
			q.Params(s.ID, gid)
		}
	})
}

func (v *Storages) AppendStorageToGroups(ctx context.Context, name string, groups ...int64) error {
	s, err := v.Get(ctx, name)
	if err != nil {
		return err
	}
	return v.db.Main().ExecContext("", ctx, func(q orm.Executor) {
		q.SQL("INSERT IGNORE INTO `storage_group` (`storage_id`, `group_id`, `created_at`, `updated_at`) " +
			"VALUES (?, ?, now(), now());")

		for _, gid := range groups {
			q.Params(s.ID, gid)
		}
	})
}

var storeNameRex = regexp.MustCompile(`^[0-9a-z\-]+$`)

func Validate(name string) bool {
	name = strings.ToLower(name)
	switch name {
	case "admin", "ui":
		return false
	default:
		return storeNameRex.MatchString(name)
	}
}
