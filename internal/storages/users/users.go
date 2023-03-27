package users

import (
	"context"
	"fmt"

	"github.com/deweppro/go-sdk/orm"
	"golang.org/x/crypto/bcrypt"
)

func (v *Users) GetUserByLogin(ctx context.Context, login string) (User, error) {
	u, ok := v.cache.GetByLogin(login)
	if ok {
		return u, nil
	}
	u, t, err := v.findUserByLogin(ctx, login)
	if err != nil {
		return User{}, err
	}
	v.cache.Setup(u, t...)
	return u, nil
}

func (v *Users) findUserByLogin(ctx context.Context, login string) (User, []string, error) {
	var (
		uid             uint64
		ulogin, upasswd string
	)
	err := v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `id`, `login`, `passwd` FROM `users` WHERE `login` = ? AND `lock` = 0 LIMIT 1;", login)
		q.Bind(func(bind orm.Scanner) error {
			return bind.Scan(&uid, &ulogin, &upasswd)
		})
	})
	if err != nil {
		return User{}, nil, err
	}
	if uid == 0 {
		return User{}, nil, fmt.Errorf("user not found")
	}

	u := User{
		ID:     uid,
		Login:  login,
		Passwd: []byte(upasswd),
		Groups: make(map[uint64]struct{}, 0),
	}

	err = v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `group_id` FROM `user_group` WHERE `user_id` = ?;", uid)
		q.Bind(func(bind orm.Scanner) error {
			var gid uint64
			if err = bind.Scan(&gid); err != nil {
				return err
			}
			u.Groups[gid] = struct{}{}
			return nil
		})
	})
	if err != nil {
		return User{}, nil, err
	}

	tokens := make([]string, 0, 10)
	err = v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `token` FROM `user_token` WHERE `user_id` = ?;", uid)
		q.Bind(func(bind orm.Scanner) error {
			var token string
			if err = bind.Scan(&token); err != nil {
				return err
			}
			tokens = append(tokens, token)
			return nil
		})
	})
	if err != nil {
		return User{}, nil, err
	}

	return u, tokens, nil
}

func (v *Users) ValidateUserPasswd(ctx context.Context, login, passwd string) bool {
	u, err := v.GetUserByLogin(ctx, login)
	if err != nil {
		return false
	}

	bp := []byte(passwd)
	err = bcrypt.CompareHashAndPassword(u.Passwd, bp)

	return err == nil
}

func (v *Users) HasUserInGroup(ctx context.Context, login string, groups ...uint64) bool {
	u, err := v.GetUserByLogin(ctx, login)
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
	u, err := v.GetUserByLogin(ctx, login)
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
	u, err := v.GetUserByLogin(ctx, login)
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
