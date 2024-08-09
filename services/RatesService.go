package services

import (
	"fmt"
	"sync"
	"time"
	"uacc-backend/integrations"
)

type CurrencyRate struct {
	From      string
	To        string
	Rate      float64
	Timestamp time.Time
	Source    string
}

// CurrencyRateLookup is a map where each key is a string, and the value is another map.
type CurrencyRateLookup map[string]map[string]CurrencyRate

type RatesService interface {
	GetRates() (CurrencyRateLookup, error)
}

type RatesServiceImpl struct {
	agents []integrations.ProxyAgent
}

func NewRatesService(agents []integrations.ProxyAgent) *RatesServiceImpl {
	return &RatesServiceImpl{agents}
}

func (service *RatesServiceImpl) GetRates() (CurrencyRateLookup, error) {
	var wg sync.WaitGroup
	resultsChannel := make(chan integrations.RateResponse, len(service.agents))

	for _, agent := range service.agents {
		wg.Add(1)
		go func(agent integrations.ProxyAgent) {
			defer wg.Done()
			rates, err := agent.GetRates()
			if err != nil {
				fmt.Println("Error getting rates:", err)
				return
			}
			resultsChannel <- rates
		}(agent)
	}

	go func() {
		wg.Wait()
		close(resultsChannel)
	}()

	result := make(CurrencyRateLookup)
	for rates := range resultsChannel {
		_, hasRate := result[rates.Base]
		if !hasRate {
			result[rates.Base] = make(map[string]CurrencyRate)
		}
		base := result[rates.Base]

		for key, rate := range rates.Rates {
			exchange := CurrencyRate{}
			exchange.Source = rates.Source
			exchange.From = rates.Base
			exchange.To = key
			exchange.Rate = rate
			exchange.Timestamp = rates.Timestamp
			base[key] = exchange
		}
	}

	return result, nil
}
