package users

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"time"

	"github.com/arwos/artifactory/internal/proxy/db"
	"github.com/deweppro/go-sdk/app"
	"github.com/deweppro/go-sdk/orm"
	"github.com/deweppro/go-sdk/random"
	"github.com/deweppro/go-sdk/routine"
	"github.com/deweppro/goppy/plugins"
	"golang.org/x/crypto/bcrypt"
)

var Plugin = plugins.Plugin{
	Inject: NewUsers,
}

type (
	Users struct {
		db db.DB

		cache *Cache
	}

	User struct {
		ID       int64
		Login    string
		Passwd   []byte
		TempKey  []byte
		TempHash []byte
		Groups   map[int64]struct{}
	}
)

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

func (v *Users) Get(ctx context.Context, login string) (User, error) {
	u, ok := v.cache.Get(login)
	if ok {
		return u, nil
	}
	u, err := v.reloadUserFromDB(ctx, login)
	if err != nil {
		return User{}, err
	}
	v.cache.Set(u)
	return u, nil
}

func (v *Users) reloadUserFromDB(ctx context.Context, login string) (User, error) {
	var (
		uid             int64
		ulogin, upasswd string
	)
	err := v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `id`, `login`, `passwd` FROM `users` WHERE `login` = ? AND `lock` = 0 LIMIT 1;", login)
		q.Bind(func(bind orm.Scanner) error {
			return bind.Scan(&uid, &ulogin, &upasswd)
		})
	})
	if err != nil {
		return User{}, err
	}
	if uid == 0 {
		return User{}, fmt.Errorf("user not found")
	}

	u := User{
		ID:       uid,
		Login:    login,
		Passwd:   []byte(upasswd),
		TempKey:  random.Bytes(64),
		TempHash: nil,
		Groups:   make(map[int64]struct{}, 0),
	}

	err = v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `group_id` FROM `user_group` WHERE `user_id` = ? LIMIT 1;", uid)
		q.Bind(func(bind orm.Scanner) error {
			var gid int64
			if err = bind.Scan(&gid); err != nil {
				return err
			}
			u.Groups[gid] = struct{}{}
			return nil
		})
	})
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (v *Users) ValidateUserPasswd(ctx context.Context, login, passwd string) bool {
	u, err := v.Get(ctx, login)
	if err != nil {
		return false
	}

	bp := []byte(passwd)

	mac := hmac.New(sha512.New, u.TempKey)
	mac.Write(bp)
	hash := mac.Sum(nil)

	if len(u.TempHash) > 0 && hmac.Equal(u.TempHash, hash) {
		return true
	}

	err = bcrypt.CompareHashAndPassword(u.Passwd, bp)
	if err != nil {
		return false
	}

	u.TempHash = hash
	v.cache.Set(u)

	return true
}

func (v *Users) HasUserInGroup(ctx context.Context, login string, groups ...int64) bool {
	u, err := v.Get(ctx, login)
	if err != nil {
		return false
	}

	for _, group := range groups {
		if _, ok := u.Groups[group]; ok {
			return true
		}
	}

	return false
}

func (v *Users) CreateUser(ctx context.Context, login, passwd string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return v.db.Main().ExecContext("", ctx, func(q orm.Executor) {
		q.SQL("INSERT INTO `users` (`login`, `passwd`, `acl`, `lock`, `created_at`, `updated_at`) "+
			"VALUES (?, ?, '', '0', now(), now());", login, string(bytes))
		q.Bind(func(result orm.Result) error {
			if result.RowsAffected == 0 {
				return fmt.Errorf("failed to create user")
			}
			return nil
		})
	})
}

func (v *Users) CreateGroup(ctx context.Context, name string) error {
	return v.db.Main().ExecContext("", ctx, func(q orm.Executor) {
		q.SQL("INSERT INTO `group` (`name`, `created_at`, `updated_at`) "+
			"VALUES (?, now(), now());", name)
		q.Bind(func(result orm.Result) error {
			if result.RowsAffected == 0 {
				return fmt.Errorf("failed to create group")
			}
			return nil
		})
	})
}

func (v *Users) ListGroup(ctx context.Context) (map[int64]string, error) {
	result := make(map[int64]string)
	err := v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `id`, `name` FROM `group`")
		q.Bind(func(bind orm.Scanner) error {
			var (
				id   int64
				name string
			)
			if err := bind.Scan(&id, &name); err != nil {
				return err
			}
			result[id] = name
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *Users) DeleteUserFromGroups(ctx context.Context, login string, groups ...int64) error {
	u, err := v.Get(ctx, login)
	if err != nil {
		return err
	}

	return v.db.Main().ExecContext("", ctx, func(q orm.Executor) {
		q.SQL("DELETE FROM `user_group` WHERE `user_id` = ? AND `group_id` = ?;")

		for _, gid := range groups {
			q.Params(u.ID, gid)
		}
	})
}

func (v *Users) AppendUserToGroups(ctx context.Context, login string, groups ...int64) error {
	u, err := v.Get(ctx, login)
	if err != nil {
		return err
	}
	return v.db.Main().ExecContext("", ctx, func(q orm.Executor) {
		q.SQL("INSERT IGNORE INTO `user_group` (`user_id`, `group_id`, `created_at`, `updated_at`) " +
			"VALUES (?, ?, now(), now());")

		for _, gid := range groups {
			q.Params(u.ID, gid)
		}
	})
}
