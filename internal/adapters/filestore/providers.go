package filestore

import (
	"fmt"

	"go.arwos.org/objex/internal/adapters/filestore/common"
	"go.arwos.org/objex/internal/adapters/filestore/ftp"
	"go.arwos.org/objex/internal/adapters/filestore/local"
	"go.osspkg.com/goppy/plugins"
	"go.osspkg.com/goppy/xc"
)

var Plugin = plugins.Plugin{
	Config: &common.Config{},
	Inject: New,
}

type (
	object struct {
		conf *common.Config
		list map[string]common.TProvider
	}

	Providers interface {
		ByCode(c string) (common.TProvider, error)
		Codes() []string
	}
)

func New(c *common.Config) Providers {
	return &object{
		conf: c,
		list: make(map[string]common.TProvider, len(c.Items)),
	}
}

func (v *object) Up(ctx xc.Context) error {
	for _, c := range v.conf.Items {
		switch c.Type {
		case common.TYPE_LOCAL:
			v.list[c.Code] = local.New(ctx.Context(), c)
		case common.TYPE_FTP:
			v.list[c.Code] = ftp.New(ctx.Context(), c)
		default:
			return fmt.Errorf("unknown provider type: %s", c.Type)
		}
	}
	for code, provider := range v.list {
		if err := provider.Check(); err != nil {
			return fmt.Errorf("provider check [%s]: %w", code, err)
		}
	}
	return nil
}

func (v *object) Down() error {
	return nil
}

func (v *object) ByCode(c string) (common.TProvider, error) {
	p, ok := v.list[c]
	if ok {
		return p, nil
	}
	return nil, fmt.Errorf("provider with code [%s] not found", c)
}

func (v *object) Codes() []string {
	result := make([]string, 0, len(v.list))
	for c := range v.list {
		result = append(result, c)
	}
	return result
}
