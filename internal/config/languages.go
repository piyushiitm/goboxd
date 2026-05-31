package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Language struct {
	Extension string   `yaml:"extension"`
	Compile   []string `yaml:"compile"`
	Run       []string `yaml:"run"`
}

func LoadLanguages() (map[string]Language, error) {

	data, err := os.ReadFile("config/languages.yaml")

	if err != nil {
		return nil, err
	}

	var languages map[string]Language

	err = yaml.Unmarshal(data, &languages)

	if err != nil {
		return nil, err
	}

	return languages, nil
}
