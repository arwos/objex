package users

import (
	"context"
	"fmt"

	"github.com/deweppro/go-sdk/orm"
	"github.com/deweppro/go-sdk/random"
)

func (v *Users) GetUserByToken(ctx context.Context, token string) (User, error) {
	u, ok := v.cache.GetByToken(token)
	if ok {
		return u, nil
	}
	login, err := v.findLoginByToken(ctx, token)
	if err != nil {
		return User{}, err
	}
	return v.GetUserByLogin(ctx, login)
}

func (v *Users) findLoginByToken(ctx context.Context, token string) (string, error) {
	login := ""
	err := v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT u.`login` FROM `users` AS u "+
			"JOIN `user_token` AS t ON u.`id` = t.`user_id` "+
			"WHERE t.`token` = ? AND u.`lock` = 0 LIMIT 1;", token)
		q.Bind(func(bind orm.Scanner) error {
			return bind.Scan(&login)
		})
	})
	if err != nil {
		return "", err
	}
	if len(login) == 0 {
		return "", fmt.Errorf("token not found")
	}
	return login, nil
}

func (v *Users) CreateToken(ctx context.Context, login string) (string, error) {
	u, err := v.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	token := random.String(64)
	err = v.db.Main().ExecContext("", ctx, func(q orm.Executor) {
		q.SQL("INSERT INTO `user_token` (`token`, `user_id`, `created_at`) "+
			"VALUES (?, ?, now());", token, u.ID)
		q.Bind(func(result orm.Result) error {
			if result.RowsAffected == 0 {
				return fmt.Errorf("failed to create token")
			}
			return nil
		})
	})
	if err != nil {
		return "", err
	}
	v.cache.SetToken(u.ID, token)
	return token, nil
}

func (v *Users) DeleteToken(ctx context.Context, login string, token string) error {
	u, err := v.GetUserByLogin(ctx, login)
	if err != nil {
		return err
	}
	err = v.db.Main().ExecContext("", ctx, func(q orm.Executor) {
		q.SQL("DELETE FROM `user_token` WHERE `user_id` = ? AND `token` = ?;", u.ID, token)
		q.Bind(func(result orm.Result) error {
			if result.RowsAffected == 0 {
				return fmt.Errorf("failed to delete token")
			}
			return nil
		})
	})
	if err != nil {
		return err
	}
	v.cache.DeleteToken(u.ID, token)
	return nil
}

func (v *Users) ListToken(ctx context.Context, uid uint64) (Tokens, error) {
	tokens := make(Tokens, 0, 10)
	err := v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT `id`, `token`, `created_at` FROM `user_token` WHERE `user_id` = ?;", uid)
		q.Bind(func(bind orm.Scanner) error {
			token := Token{}
			if err := bind.Scan(&token.ID, &token.Token, &token.CreatedAt); err != nil {
				return err
			}
			tokens = append(tokens, token)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
