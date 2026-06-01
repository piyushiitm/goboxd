package executor

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/piyushiitm/goboxd/internal/config"
	"github.com/piyushiitm/goboxd/internal/models"
	"github.com/piyushiitm/goboxd/internal/validator"
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

func ResolveLimits(
	defaults config.Limits,
	override *models.Limits,
) config.Limits {

	result := defaults

	if override == nil {
		return result
	}

	if override.WallTimeS > 0 {
		result.WallTimeS = override.WallTimeS
	}

	if override.MemoryKB > 0 {
		result.MemoryKB = override.MemoryKB
	}

	if override.MaxProcesses > 0 {
		result.MaxProcesses = override.MaxProcesses
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
		return models.RunResponse{}, validator.ErrUnknownLanguage
	}
	effectiveRunLimits := ResolveLimits(
		languageConfig.DefaultRunLimits,
		nil,
	)
	if req.Run != nil {
		effectiveRunLimits = ResolveLimits(
			languageConfig.DefaultRunLimits,
			req.Run.Limits,
		)
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

	buildResult := &models.BuildResult{
		Status: "ok",
		Stdout: "",
		Stderr: "",
	}

	buildStart := time.Now()

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
				Build: &models.BuildResult{
					Status:     "compile_error",
					Stdout:     "",
					Stderr:     errorMessage,
					DurationMS: time.Since(buildStart).Milliseconds(),
				},
			}, nil
		}
	}

	buildResult.DurationMS = time.Since(buildStart).Milliseconds()

	results := []models.TestResult{}
	overallStatus := "accepted"

	for _, test := range req.Tests {
		testStart := time.Now()
		runCommand := ReplacePlaceholders(
			languageConfig.Run,
			sourceFile,
			artifactFile,
		)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(effectiveRunLimits.WallTimeS)*time.Second,
		)
		defer cancel()

		cmd := exec.CommandContext(
			ctx,
			runCommand[0],
			runCommand[1:]...,
		)

		cmd.Stdin = strings.NewReader(test.Stdin)

		output, err := cmd.CombinedOutput()

		if ctx.Err() == context.DeadlineExceeded {

			results = append(results, models.TestResult{
				Status:       "time_limit_exceeded",
				Stdout:       "",
				Stderr:       "",
				DurationMS:   time.Since(testStart).Milliseconds(),
				MemoryPeakKB: 0,
			})

			overallStatus = "time_limit_exceeded"

			continue
		}

		if err != nil {
			return models.RunResponse{
				Status: "runtime_error",
				Build:  buildResult,
			}, nil
		}

		actual := strings.TrimSpace(string(output))
		expected := strings.TrimSpace(test.ExpectedStdout)

		testStatus := "accepted"

		if actual != expected {
			testStatus = "wrong_output"
			overallStatus = "wrong_output"
		}
		duration := time.Since(testStart).Milliseconds()

		results = append(results, models.TestResult{
			Status:     testStatus,
			Stdout:     actual,
			Stderr:     "",
			DurationMS: duration,
		})
	}

	return models.RunResponse{
		Status: overallStatus,
		Build:  buildResult,
		Tests:  results,
	}, nil
}
