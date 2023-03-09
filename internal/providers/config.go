package providers

import (
	"fmt"
)

type (
	Config struct {
		Items []Item `yaml:"store_providers"`
	}

	Item struct {
		Name    string `yaml:"name"`
		Type    string `yaml:"type"`
		Code    string `yaml:"code"`
		Prefix  string `yaml:"prefix"`
		Setting string `yaml:"setting"`
	}
)

func (v *Config) Default() {
	if len(v.Items) == 0 {
		v.Items = append(v.Items, Item{
			Name:    "Local",
			Type:    TypeLocal,
			Code:    "loc",
			Prefix:  "",
			Setting: "/tmp",
		})
	}
}

func (v *Config) Validate() error {
	if len(v.Items) == 0 {
		return fmt.Errorf("providers list is empty")
	}
	for _, item := range v.Items {
		if !IsValidType(item.Type) {
			return fmt.Errorf("unknown provider type")
		}
		if len(item.Setting) == 0 {
			return fmt.Errorf("provider setting is empty")
		}
	}
	return nil
}
