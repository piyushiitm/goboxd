package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

type VersionCommand struct {
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
}

type Limits struct {
	WallTimeS    int `yaml:"wall_time_s" json:"wall_time_s"`
	MemoryKB     int `yaml:"memory_kb" json:"memory_kb"`
	MaxProcesses int `yaml:"max_processes" json:"max_processes"`
}

type Language struct {
	Name             string         `yaml:"name"`
	Extension        string         `yaml:"extension"`
	VersionCommand   VersionCommand `yaml:"version_command"`
	DefaultRunLimits Limits         `yaml:"default_run_limits"`
	Compile          []string       `yaml:"compile"`
	Run              []string       `yaml:"run"`
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
