package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

type Language struct {
	Extension string   `yaml:"extension"`
	Compile   []string `yaml:"compile"`
	Run       []string `yaml:"run"`
}

//go:embed languages.yaml
var languagesYAML []byte

func LoadLanguages() (map[string]Language, error) {

	var languages map[string]Language

	err := yaml.Unmarshal(
		languagesYAML,
		&languages,
	)

	if err != nil {
		return nil, err
	}

	return languages, nil
}
