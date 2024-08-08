package routing

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"uacc-backend/contract"
)

func SetupRates(
	data *contract.Data,
	mux *http.ServeMux,
) {
	mux.HandleFunc("/api/v4/rate/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path[len("/api/v4/rate/"):]
		parts := strings.Split(path, "/")
		if len(parts) != 2 {
			http.Error(w, "Expect path like /api/v4/rate/{from}/{to}", http.StatusBadRequest)
			return
		}
		from := parts[0]
		to := parts[1]

		result, hasRate := data.Rates[from][to]
		log.Printf("Handling %s -> %s %+v\n", from, to, result)
		if !hasRate {
			http.Error(w, "Rate not found", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/api/v5/rates/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path[len("/api/v5/rates/"):]
		to := path

		result, hasRate := data.NewRates[to]
		if !hasRate {
			http.Error(w, "Rates not found for "+to, http.StatusNotFound)
		}

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
