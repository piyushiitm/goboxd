package executor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/piyushiitm/goboxd/internal/config"
	"github.com/piyushiitm/goboxd/internal/models"
)

func ReplacePlaceholders(
	command []string,
	sourceFile string,
	artifactFile string,
) []string {

	result := make([]string, len(command))

	for i, part := range command {

		part = strings.ReplaceAll(
			part,
			"{source}",
			sourceFile,
		)

		part = strings.ReplaceAll(
			part,
			"{artifact}",
			artifactFile,
		)

		result[i] = part
	}

	return result
}

func Execute(language string, source string) (models.RunResponse, error) {
	languages, err := config.LoadLanguages()

	if err != nil {
		return models.RunResponse{}, err
	}

	languageConfig, exists := languages[language]

	if !exists {
		return models.RunResponse{}, errors.New("unsupported language")
	}

	id := uuid.New().String()
	workspace := filepath.Join("workspace", id)
	err = os.MkdirAll(workspace, 0755)
	if err != nil {
		return models.RunResponse{}, err
	}

	sourceFile := filepath.Join(
		workspace,
		"source"+languageConfig.Extension,
	)

	artifactFile := filepath.Join(
		workspace,
		"program",
	)

	err = os.WriteFile(
		sourceFile,
		[]byte(source),
		0644,
	)

	if err != nil {
		return models.RunResponse{}, err
	}

	runCommand := ReplacePlaceholders(
		languageConfig.Run,
		sourceFile,
		artifactFile,
	)

	cmd := exec.Command(
		runCommand[0],
		runCommand[1:]...,
	)

	output, err := cmd.CombinedOutput()

	fmt.Println("Created:", sourceFile)
	fmt.Println("Language:", language)

	if err != nil {
		return models.RunResponse{
			Stdout:   "",
			Stderr:   string(output),
			ExitCode: 1,
		}, nil
	}

	return models.RunResponse{
		Stdout:   string(output),
		Stderr:   "",
		ExitCode: 0,
	}, nil
}
