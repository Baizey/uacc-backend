package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	integrations2 "uacc-backend/integrations"
	services2 "uacc-backend/services"
)

type PathResponse struct {
	Source    string  `json:"source"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Rate      float64 `json:"rate"`
	Timestamp int64   `json:"timestamp"`
}

type RateResponse struct {
	From      string         `json:"from"`
	To        string         `json:"to"`
	Rate      float64        `json:"rate"`
	Timestamp int64          `json:"timestamp"`
	Path      []PathResponse `json:"path"`
}

type Data struct {
	symbols services2.SymbolsResponse
	rates   map[string]map[string]RateResponse
}

var data = &Data{}

func main() {
	log.Println("Starting up...")

	openExchangeAgent := integrations2.NewOpenExchangeProxyAgent(getOrCrash("openExchangeApiKey"))
	agents := []integrations2.ProxyAgent{openExchangeAgent}
	ratesService := services2.NewRatesService(agents)
	symbolsService := services2.NewSymbolsService(agents)

	log.Println("Getting initial data...")
	err := update(symbolsService, ratesService, data)
	if err != nil {
		log.Fatal("Error updating data:", err)
	}

	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := update(symbolsService, ratesService, data); err != nil {
					log.Println("Error updating data:", err)
				}
			}
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/v4/symbols", handleSymbols)
	mux.HandleFunc("/api/v4/rate/", handleRates)
	wrapped := validateApikey(mux)

	log.Println("Starting up server...")
	port := getOrDefault("PORT", "3001")
	log.Fatal(http.ListenAndServe(":"+port, wrapped))
}

func update(symbolsService services2.SymbolsService, ratesService services2.RatesService, data *Data) error {
	log.Println("Updating...")
	symbols, err := symbolsService.GetSymbols()
	if err != nil {
		return err
	}
	data.symbols = symbols

	rates, err := ratesService.GetRates()
	if err != nil {
		return err
	}
	tmp := make(map[string]map[string]RateResponse)
	for from := range rates {
		_, err := tmp[from]
		if !err {
			tmp[from] = make(map[string]RateResponse)
		}
		for to, fromToRate := range rates[from] {
			_, hasFromTo := tmp[from][to]
			if hasFromTo {
				item := RateResponse{}
				item.From = from
				item.To = to
				item.Timestamp = fromToRate.Timestamp.Unix()
				item.Path = make([]PathResponse, 0)
				item.Rate = fromToRate.Rate
				tmp[from][to] = item
				//log.Printf("%s -> %s %f\n", from, to, item.rate)
			}

			_, hasTo := tmp[to]
			if !hasTo {
				tmp[to] = make(map[string]RateResponse)
			}

			_, hasToFrom := tmp[to][from]
			if !hasToFrom {
				item := RateResponse{}
				item.From = to
				item.To = from
				item.Timestamp = fromToRate.Timestamp.Unix()
				item.Path = make([]PathResponse, 0)
				item.Rate = 1. / fromToRate.Rate
				tmp[to][from] = item
				//log.Printf("%s -> %s %f\n", to, from, item.rate)
			}

			for toOther, fromOtherRate := range rates[from] {
				_, hasToOther := tmp[to][toOther]
				if !hasToOther {
					item := RateResponse{}
					item.From = to
					item.To = toOther
					item.Timestamp = fromToRate.Timestamp.Unix()
					item.Path = make([]PathResponse, 0)
					item.Rate = tmp[to][from].Rate * fromOtherRate.Rate
					tmp[to][toOther] = item
				}
			}
		}
	}
	data.rates = tmp
	return nil
}

func validateApikey(next http.Handler) http.Handler {
	apikey := getOrCrash("ownApiKey")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		key := r.Header.Get("x-apikey")
		if key != apikey {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleSymbols(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result := data.symbols
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func handleRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path[len("/api/v4/rate/"):]
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		http.Error(w, "Expect path like /api/v4/rate/{from}/{to}", http.StatusBadRequest)
		return
	}
	from := parts[0]
	to := parts[1]

	result, hasRate := data.rates[from][to]
	log.Printf("Handling %s -> %s %+v\n", from, to, result)
	if !hasRate {
		http.Error(w, "Rate not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getOrCrash(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing env variable %s", key)
	}
	return value
}
func getOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Missing env variable %s using default %s", key, defaultValue)
		return defaultValue
	}
	log.Printf("Have env variable %s using %s", key, defaultValue)
	return value
}
