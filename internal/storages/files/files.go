package files

//go:generate easyjson

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/arwos/artifactory/internal/proxy/db"
	"github.com/deweppro/go-sdk/orm"
	"github.com/deweppro/goppy/plugins"
)

var Plugin = plugins.Plugin{
	Inject: NewFiles,
}

type (
	Files struct {
		db db.DB
	}

	//easyjson:json
	File struct {
		ID    int64             `json:"id"`
		Name  string            `json:"name"`
		Hash  string            `json:"sha1"`
		Props map[string]string `json:"props"`
	}
)

func NewFiles(db db.DB) *Files {
	return &Files{
		db: db,
	}
}

const querySearchByProps = "SELECT `files_id` FROM `props` WHERE `name` = ? AND `value` = ?"

func (v *Files) SearchByProps(ctx context.Context, sid int64, props url.Values) ([]File, error) {
	queries := make([]string, 0, len(props))
	params := make([]interface{}, 0, len(props)*2)

	for key := range props {
		queries = append(queries, querySearchByProps)
		params = append(params, key, props.Get(key))
	}

	ids := make([]interface{}, 0)
	err := v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL(strings.Join(queries, "\n INTERSECT \n"), params...)
		q.Bind(func(bind orm.Scanner) error {
			var id int64
			if err := bind.Scan(&id); err != nil {
				return err
			}
			ids = append(ids, id)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	result := make([]File, 0, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	query := strings.Trim(strings.Repeat("?", len(ids)), ",")
	query = fmt.Sprintf("SELECT `id`, `name`, `hash` FROM `files` WHERE `id` IN (%s);", query)

	err = v.db.Main().QueryContext("", ctx, func(q orm.Querier) {
		q.SQL(query, ids...)
		q.Bind(func(bind orm.Scanner) error {
			file := File{
				ID:    0,
				Name:  "",
				Hash:  "",
				Props: nil,
			}
			if err = bind.Scan(&file.ID, &file.Name, &file.Hash); err != nil {
				return err
			}
			result = append(result, file)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (v *Files) AddFile(ctx context.Context, sid int64, filename, hash string, props url.Values) error {
	return v.db.Main().TransactionContext("", ctx, func(v orm.Tx) {
		var id int64
		v.Exec(func(e orm.Executor) {
			e.SQL("INSERT INTO `files` (`storage_id`, `name`, `hash`, `created_at`, `updated_at`) "+
				"VALUES (?, ?, ?, now(), now());", sid, filename, hash)
			e.Bind(func(result orm.Result) error {
				if result.RowsAffected == 0 {
					return fmt.Errorf("file record error")
				}
				id = result.LastInsertId
				return nil
			})
		})

		v.Exec(func(e orm.Executor) {
			e.SQL("INSERT INTO `props` (`files_id`, `name`, `value`, `created_at`, `updated_at`) " +
				"VALUES (?, ?, ?, now(), now());")
			for key := range props {
				e.Params(id, key, props.Get(key))
			}
		})
	})
}
