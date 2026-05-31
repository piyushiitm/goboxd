package api

import (
	"encoding/json"
	"net/http"

	"github.com/piyushiitm/goboxd/internal/executor"
	"github.com/piyushiitm/goboxd/internal/models"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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
