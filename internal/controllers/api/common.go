package api

import "encoding/json"

var EMPTY = json.RawMessage(`{}`)

type (
	Config struct {
		Settings Settings `yaml:"api_settings"`
	}

	Settings struct {
		CookieName string `yaml:"cookie_name"`
		UsedHTTPS  bool   `yaml:"used_https"`
	}
)

func (v *Config) Default() {
	if len(v.Settings.CookieName) == 0 {
		v.Settings.CookieName = "sess"
	}
}
