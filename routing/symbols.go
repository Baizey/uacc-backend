package routing

import (
	"encoding/json"
	"net/http"
	"uacc-backend/contract"
)

func SetupSymbols(
	data *contract.Data,
	mux *http.ServeMux,
) {
	mux.HandleFunc("/api/v4/symbols", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		result := data.Symbols
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})
}
