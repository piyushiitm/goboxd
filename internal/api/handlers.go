package api

import (
	"encoding/json"
	"net/http"

	"github.com/piyushiitm/goboxd/internal/executor"
	"github.com/piyushiitm/goboxd/internal/models"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"build_info": map[string]string{
			"version": "0.1.0",
			"commit":  "dev",
		},
		"languages": []string{
			"cpp",
			"py3",
		},
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

	response, err := executor.Execute(
		req.Language,
		req.Source,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
