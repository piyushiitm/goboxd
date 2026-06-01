package executor

import (
	"errors"
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

func Execute(req models.RunRequest) (models.RunResponse, error) {
	language := req.Language
	source := req.Source
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

	sourceName := req.SourceFilename

	if sourceName == "" {
		sourceName = "source" + languageConfig.Extension
	}

	sourceFile := filepath.Join(
		workspace,
		sourceName,
	)

	artifactName := req.ArtifactFilename

	if artifactName == "" {
		artifactName = "program"
	}

	artifactFile := filepath.Join(
		workspace,
		artifactName,
	)

	err = os.WriteFile(
		sourceFile,
		[]byte(source),
		0644,
	)

	if err != nil {
		return models.RunResponse{}, err
	}

	if len(languageConfig.Compile) > 0 {

		compileCommand := ReplacePlaceholders(
			languageConfig.Compile,
			sourceFile,
			artifactFile,
		)

		cmd := exec.Command(
			compileCommand[0],
			compileCommand[1:]...,
		)

		output, err := cmd.CombinedOutput()

		if err != nil {

			errorMessage := string(output)

			if errorMessage == "" {
				errorMessage = err.Error()
			}

			errorMessage = strings.ReplaceAll(
				errorMessage,
				sourceFile,
				"source"+languageConfig.Extension,
			)

			errorMessage = strings.ReplaceAll(
				errorMessage,
				filepath.Dir(sourceFile)+"/",
				"",
			)

			return models.RunResponse{
				Status: "compile_error",
			}, nil
		}
	}

	results := []models.TestResult{}
	overallStatus := "accepted"

	for _, test := range req.Tests {

		runCommand := ReplacePlaceholders(
			languageConfig.Run,
			sourceFile,
			artifactFile,
		)

		cmd := exec.Command(
			runCommand[0],
			runCommand[1:]...,
		)

		cmd.Stdin = strings.NewReader(test.Stdin)

		output, err := cmd.CombinedOutput()

		if err != nil {
			return models.RunResponse{
				Status: "runtime_error",
			}, nil
		}

		actual := strings.TrimSpace(string(output))
		expected := strings.TrimSpace(test.ExpectedStdout)

		testStatus := "accepted"

		if actual != expected {
			testStatus = "wrong_output"
			overallStatus = "wrong_output"
		}

		results = append(results, models.TestResult{
			Status: testStatus,
			Stdout: actual,
			Stderr: "",
		})
	}

	return models.RunResponse{
		Status: overallStatus,
		Tests:  results,
	}, nil
}
