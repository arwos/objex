package providers

import (
	"fmt"

	"github.com/deweppro/go-sdk/app"
)

type (
	providers struct {
		conf *Config
		list map[string]Provider
	}

	Providers interface {
		GetByCode(c string) (Provider, error)
		ListCodes() []string
	}
)

func New(c *Config) (*providers, Providers) {
	p := &providers{
		conf: c,
		list: make(map[string]Provider, len(c.Items)),
	}
	return p, p
}

func (v *providers) Up(ctx app.Context) error {
	for _, c := range v.conf.Items {
		switch c.Type {
		case TypeLocal:
			v.list[c.Code] = NewLocalProvider(ctx.Context(), c)
		case TypeFTP:
			v.list[c.Code] = NewFTPProvider(ctx.Context(), c)
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

func (v *providers) Down() error {
	return nil
}

func (v *providers) GetByCode(c string) (Provider, error) {
	p, ok := v.list[c]
	if ok {
		return p, nil
	}
	return nil, fmt.Errorf("provider with code [%s] not found", c)
}

func (v *providers) ListCodes() []string {
	result := make([]string, 0, len(v.list))
	for c := range v.list {
		result = append(result, c)
	}
	return result
}
