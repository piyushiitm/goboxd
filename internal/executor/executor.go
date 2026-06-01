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
	"github.com/piyushiitm/goboxd/internal/sandbox"
	"github.com/piyushiitm/goboxd/internal/validator"
)

const (
	MaxStdoutBytes = 64 * 1024
	MaxStderrBytes = 64 * 1024
	MaxOutputBytes = 128 * 1024
)

func ReplacePlaceholders(
	command []string,
	sourceFile string,
	artifactFile string,
	workspace string,
	buildFlags []string,
	runFlags []string,
) []string {

	result := []string{}

	for _, part := range command {

		if part == "{build_flags}" {
			result = append(result, buildFlags...)
			continue
		}

		if part == "{run_flags}" {
			result = append(result, runFlags...)
			continue
		}

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

		part = strings.ReplaceAll(
			part,
			"{workspace}",
			workspace,
		)

		result = append(result, part)
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

type LimitedBuffer struct {
	Data  []byte
	Limit int
}

func (b *LimitedBuffer) Write(p []byte) (int, error) {

	remaining := b.Limit - len(b.Data)

	if remaining > 0 {

		if len(p) > remaining {
			b.Data = append(
				b.Data,
				p[:remaining]...,
			)
		} else {
			b.Data = append(
				b.Data,
				p...,
			)
		}
	}

	return len(p), nil
}

func NotExecutedTests(
	tests []models.TestCase,
) []models.TestResult {

	results := make(
		[]models.TestResult,
		len(tests),
	)

	for i := range tests {
		results[i] = models.TestResult{
			Status: "not_executed",
		}
	}

	return results
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

	if req.Build != nil {

		if !validator.ValidateFlags(
			req.Build.Flags,
			languageConfig.AllowedBuildFlags,
		) {
			return models.RunResponse{},
				validator.ErrInvalidBuildFlag
		}
	}

	if req.Run != nil {

		if !validator.ValidateFlags(
			req.Run.Flags,
			languageConfig.AllowedRunFlags,
		) {
			return models.RunResponse{},
				validator.ErrInvalidRunFlag
		}
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

	defer os.RemoveAll(workspace)

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

		buildFlags := []string{}

		if req.Build != nil {
			buildFlags = req.Build.Flags
		}

		compileCommand := ReplacePlaceholders(
			languageConfig.Compile,
			sourceFile,
			artifactFile,
			workspace,
			buildFlags,
			nil,
		)

		cmd := exec.Command(
			compileCommand[0],
			compileCommand[1:]...,
		)

		var stdout LimitedBuffer
		var stderr LimitedBuffer

		stdout.Limit = MaxStdoutBytes
		stderr.Limit = MaxStderrBytes

		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {

			errorMessage := string(stderr.Data)

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
				Status: "build_failed",
				Build: &models.BuildResult{
					Status:     "failed",
					Stdout:     "",
					Stderr:     errorMessage,
					DurationMS: time.Since(buildStart).Milliseconds(),
				},
				Tests: NotExecutedTests(req.Tests),
			}, nil
		}
	}

	buildResult.DurationMS = time.Since(buildStart).Milliseconds()

	results := []models.TestResult{}
	overallStatus := "accepted"

	for _, test := range req.Tests {
		testStart := time.Now()
		runFlags := []string{}

		if req.Run != nil {
			runFlags = req.Run.Flags
		}

		runCommand := ReplacePlaceholders(
			languageConfig.Run,
			sourceFile,
			artifactFile,
			workspace,
			nil,
			runFlags,
		)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(effectiveRunLimits.WallTimeS)*time.Second,
		)
		defer cancel()

		runner := sandbox.NativeRunner{}

		cmd := runner.Command(
			runCommand,
			workspace,
			effectiveRunLimits.WallTimeS,
		)

		cmd.Stdin = strings.NewReader(test.Stdin)

		var stdout LimitedBuffer
		var stderr LimitedBuffer

		stdout.Limit = MaxStdoutBytes
		stderr.Limit = MaxStderrBytes

		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if ctx.Err() == context.DeadlineExceeded {

			results = append(results, models.TestResult{
				Status:       "time_exceeded",
				Stdout:       "",
				Stderr:       "",
				DurationMS:   time.Since(testStart).Milliseconds(),
				MemoryPeakKB: 0,
			})

			overallStatus = "time_exceeded"

			continue
		}

		if err != nil {
			return models.RunResponse{
				Status: "runtime_error",
				Build:  buildResult,
				Tests: []models.TestResult{
					{
						Status: "runtime_error",
						Stderr: string(stderr.Data),
					},
				},
			}, nil
		}

		actual := strings.TrimSpace(
			string(stdout.Data),
		)
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
