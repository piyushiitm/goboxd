package api

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/piyushiitm/goboxd/internal/config"
	"github.com/piyushiitm/goboxd/internal/executor"
	"github.com/piyushiitm/goboxd/internal/models"
	"github.com/piyushiitm/goboxd/internal/validator"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	languages, err := config.LoadLanguages()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := map[string]interface{}{}
	status := "ok"
	for id, lang := range languages {
		cmd := exec.Command(
			lang.VersionCommand.Command,
			lang.VersionCommand.Args...,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {

			results[id] = map[string]interface{}{
				"ok":    false,
				"error": strings.TrimSpace(string(output)),
			}

			status = "degraded"
			continue
		}

		results[id] = map[string]interface{}{
			"ok":      true,
			"version": strings.TrimSpace(string(output)),
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if status == "degraded" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    status,
		"languages": results,
	})
}

func Info(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	languages, err := config.LoadLanguages()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	languageList := []map[string]interface{}{}

	for id, lang := range languages {

		version := ""

		cmd := exec.Command(
			lang.VersionCommand.Command,
			lang.VersionCommand.Args...,
		)

		output, err := cmd.CombinedOutput()

		if err == nil {
			version = strings.TrimSpace(string(output))
		}

		languageList = append(
			languageList,
			map[string]interface{}{
				"id":                 id,
				"name":               lang.Name,
				"version":            version,
				"default_run_limits": lang.DefaultRunLimits,
			},
		)
	}

	response := map[string]interface{}{
		"build_info": map[string]string{
			"version":    "0.1.0",
			"commit":     "dev",
			"go_version": runtime.Version(),
		},
		"languages": languageList,
		"limits": map[string]int{
			"max_source_bytes":    262144,
			"max_tests":           50,
			"max_concurrent_jobs": 16,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func writeError(
	w http.ResponseWriter,
	code string,
	message string,
	status int,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(
		models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:    code,
				Message: message,
			},
		},
	)
}

func Run(w http.ResponseWriter, r *http.Request) {

	var req models.RunRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		writeError(
			w,
			"invalid_json",
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	err = validator.Validate(req)

	if err != nil {

		switch err {

		case validator.ErrMissingLanguage:
			writeError(
				w,
				"missing_language",
				err.Error(),
				http.StatusBadRequest,
			)

		case validator.ErrMissingSource:
			writeError(
				w,
				"missing_source",
				err.Error(),
				http.StatusBadRequest,
			)

		case validator.ErrMissingTests:
			writeError(
				w,
				"missing_tests",
				err.Error(),
				http.StatusBadRequest,
			)

		case validator.ErrInvalidSourceFilename:
			writeError(
				w,
				"invalid_filename",
				err.Error(),
				http.StatusBadRequest,
			)

		case validator.ErrInvalidArtifactFilename:
			writeError(
				w,
				"invalid_filename",
				err.Error(),
				http.StatusBadRequest,
			)

		default:
			writeError(
				w,
				"validation_error",
				err.Error(),
				http.StatusBadRequest,
			)
		}

		return
	}

	response, err := executor.Execute(req)
	if err != nil {

		switch err {

		case validator.ErrUnknownLanguage:
			writeError(
				w,
				"unknown_language",
				err.Error(),
				http.StatusBadRequest,
			)

		case validator.ErrInvalidBuildFlag:
			writeError(
				w,
				"invalid_build_flag",
				err.Error(),
				http.StatusBadRequest,
			)

		case validator.ErrInvalidRunFlag:
			writeError(
				w,
				"invalid_run_flag",
				err.Error(),
				http.StatusBadRequest,
			)

		default:
			writeError(
				w,
				"internal_error",
				err.Error(),
				http.StatusInternalServerError,
			)
		}

		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
