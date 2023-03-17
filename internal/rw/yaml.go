package rw

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadYaml(filename string, data interface{}) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, data)
}

func WriteYaml(filename string, data interface{}) error {
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0644)
}
