package handler

import (
	"encoding/json"
	"net/http"

	"github.com/michaelzhan1/recent-max/internal/value/generate"
)

func ResetGeneratorHandler(g *generate.Generator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type request struct {
			Value float64 `json:"value"`
			Mu    float64 `json:"mu"`
			Sigma float64 `json:"sigma"`
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		g.Reset(req.Value, req.Mu, req.Sigma)
		w.WriteHeader(http.StatusOK)
	}
}
