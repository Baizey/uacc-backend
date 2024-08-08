package main

import (
	"log"
	"net/http"
	"time"
	"uacc-backend/contract"
	"uacc-backend/integrations"
	"uacc-backend/routing"
	"uacc-backend/services"
	"uacc-backend/util"
)

func main() {
	data := &contract.Data{}
	log.Println("Starting up...")

	agents := []integrations.ProxyAgent{
		integrations.NewOpenExchangeProxyAgent(util.GetOrCrash("openExchangeApiKey")),
	}
	ratesService := services.NewRatesService(agents)
	symbolsService := services.NewSymbolsService(agents)

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
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	routing.SetupSymbols(data, mux)
	routing.SetupRates(data, mux)
	routing.SetupLocalizations(mux)
	wrapped := routing.SetupMiddleware(mux)

	log.Println("Starting up server...")
	port := util.GetOrDefault("PORT", "3001")
	log.Fatal(http.ListenAndServe(":"+port, wrapped))
}

func update(symbolsService services.SymbolsService, ratesService services.RatesService, data *contract.Data) error {
	log.Println("Updating...")
	symbols, err := symbolsService.GetSymbols()
	if err != nil {
		return err
	}
	data.Symbols = symbols

	rates, err := ratesService.GetRates()
	if err != nil {
		return err
	}
	tmpRates := make(map[string]map[string]contract.RateResponse)
	for from := range rates {
		_, err := tmpRates[from]
		if !err {
			tmpRates[from] = make(map[string]contract.RateResponse)
		}
		for to, fromToRate := range rates[from] {
			_, hasFromTo := tmpRates[from][to]
			if hasFromTo {
				item := contract.RateResponse{}
				item.From = from
				item.To = to
				item.Timestamp = fromToRate.Timestamp.Unix()
				item.Rate = fromToRate.Rate
				item.Path = []contract.PathResponse{
					{
						From:      item.From,
						To:        item.To,
						Rate:      item.Rate,
						Source:    fromToRate.Source,
						Timestamp: item.Timestamp,
					},
				}
				tmpRates[from][to] = item
			}

			_, hasTo := tmpRates[to]
			if !hasTo {
				tmpRates[to] = make(map[string]contract.RateResponse)
			}

			_, hasToFrom := tmpRates[to][from]
			if !hasToFrom {
				item := contract.RateResponse{}
				item.From = to
				item.To = from
				item.Timestamp = fromToRate.Timestamp.Unix()
				item.Rate = 1. / fromToRate.Rate
				item.Path = []contract.PathResponse{
					{
						From:      item.To,
						To:        item.From,
						Rate:      item.Rate,
						Source:    fromToRate.Source,
						Timestamp: item.Timestamp,
					},
				}
				tmpRates[to][from] = item
			}

			for toOther, fromOtherRate := range rates[from] {
				_, hasToOther := tmpRates[to][toOther]
				if !hasToOther {
					item := contract.RateResponse{}
					item.From = to
					item.To = toOther
					item.Timestamp = fromToRate.Timestamp.Unix()
					item.Path = make([]contract.PathResponse, 0)
					item.Rate = tmpRates[to][from].Rate * fromOtherRate.Rate
					item.Path = []contract.PathResponse{
						{
							From:      to,
							To:        from,
							Rate:      tmpRates[to][from].Rate,
							Source:    fromToRate.Source,
							Timestamp: item.Timestamp,
						},
						{
							From:      from,
							To:        toOther,
							Rate:      fromOtherRate.Rate,
							Source:    fromOtherRate.Source,
							Timestamp: item.Timestamp,
						},
					}
					tmpRates[to][toOther] = item
				}
			}
		}
	}
	data.Rates = tmpRates

	tmpNewRates := make(map[string]contract.RatesResponse)
	for to := range tmpRates {
		rateResponses := make([]contract.RateResponse, len(tmpRates))
		i := 0
		for from := range tmpRates {
			rateResponses[i] = tmpRates[from][to]
			i++
		}
		resp := contract.RatesResponse{}
		resp.Rates = rateResponses
		tmpNewRates[to] = resp
	}
	data.NewRates = tmpNewRates

	return nil
}
