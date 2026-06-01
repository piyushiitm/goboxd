package api

import (
	"encoding/json"
	"net/http"
	"os/exec"
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

	languageList := []map[string]string{}

	for id := range languages {
		languageList = append(
			languageList,
			map[string]string{
				"id": id,
			},
		)
	}

	response := map[string]interface{}{
		"build_info": map[string]string{
			"version": "0.1.0",
			"commit":  "dev",
		},
		"languages": languageList,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Run(w http.ResponseWriter, r *http.Request) {

	var req models.RunRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = validator.Validate(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := executor.Execute(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
