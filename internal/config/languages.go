package config

import (
	_ "embed"
	"os/exec"

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

	AllowedBuildFlags []string `yaml:"allowed_build_flags"`
	AllowedRunFlags   []string `yaml:"allowed_run_flags"`

	Compile []string `yaml:"compile"`
	Run     []string `yaml:"run"`
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

	for id, lang := range languages {

		path, err := exec.LookPath(
			lang.VersionCommand.Command,
		)

		if err == nil {
			lang.VersionCommand.Command = path
		}

		if len(lang.Compile) > 0 {

			path, err := exec.LookPath(
				lang.Compile[0],
			)

			if err == nil {
				lang.Compile[0] = path
			}
		}

		if len(lang.Run) > 0 {

			path, err := exec.LookPath(
				lang.Run[0],
			)

			if err == nil {
				lang.Run[0] = path
			}
		}

		languages[id] = lang
	}

	return languages, nil
}
