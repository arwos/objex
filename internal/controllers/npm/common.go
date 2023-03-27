package npm

import (
	"fmt"

	"github.com/deweppro/goppy/plugins/web"
)

const (
	registryURI = "https://registry.npmjs.org"

	registry      = "/npm"
	registryFiles = "/@npm"
)

type (
	Config struct {
		Packages Packages `yaml:"npm_packages"`
	}

	Packages struct {
		SSL        bool   `yaml:"ssl"`
		ProxyCache string `yaml:"proxy_cache"`
	}
)

func (v *Config) Default() {
	if len(v.Packages.ProxyCache) == 0 {
		v.Packages.ProxyCache = "/tmp/npm"
	}
}

func (v *Config) URISchema() string {
	if v.Packages.SSL {
		return "https://"
	}
	return "http://"
}

func dump(c web.Context) {
	fmt.Println(c.URL().Path)
	b := make([]byte, 0)
	err := c.BindBytes(&b)
	fmt.Println(err, string(b))
}
