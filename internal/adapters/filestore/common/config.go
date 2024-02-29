package common

import "fmt"

type (
	Config struct {
		Items []ConfigItem `yaml:"storages"`
	}

	ConfigItem struct {
		Name    string `yaml:"name"`
		Type    string `yaml:"type"`
		Code    string `yaml:"code"`
		Setting string `yaml:"setting"`
	}
)

func (v *Config) Default() {
	if len(v.Items) == 0 {
		v.Items = append(v.Items, ConfigItem{
			Name:    "Local",
			Type:    TYPE_LOCAL,
			Code:    "loc",
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
