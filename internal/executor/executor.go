package executor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/piyushiitm/goboxd/internal/models"
)

type Language struct {
	Extension string
	Compile   []string
	Run       []string
}

var Languages = map[string]Language{
	"python": {
		Extension: ".py",
		Run:       []string{"python3"},
	},

	"cpp": {
		Extension: ".cpp",
		Compile:   []string{"g++"},
		Run:       []string{"./program"},
	},
}

func Execute(language string, source string) (models.RunResponse, error) {

	languageConfig, exists := Languages[language]

	if !exists {
		return models.RunResponse{}, errors.New("unsupported language")
	}

	id := uuid.New().String()
	workspace := filepath.Join("workspace", id)
	err := os.MkdirAll(workspace, 0755)
	if err != nil {
		return models.RunResponse{}, err
	}

	sourceFile := filepath.Join(
		workspace,
		"source"+languageConfig.Extension,
	)

	err = os.WriteFile(
		sourceFile,
		[]byte(source),
		0644,
	)

	if err != nil {
		return models.RunResponse{}, err
	}

	cmd := exec.Command(
		languageConfig.Run[0],
		sourceFile,
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
