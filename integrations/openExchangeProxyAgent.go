package integrations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type IncomingRatesResponse struct {
	Disclaimer string             `json:"disclaimer"`
	Licence    string             `json:"licence"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

type IncomingSymbolsResponse map[string]string

// OpenExchangeProxyAgent implements the proxy agent interface
type OpenExchangeProxyAgent struct {
	apiKey string
}

func NewOpenExchangeProxyAgent(apikey string) *OpenExchangeProxyAgent {
	return &OpenExchangeProxyAgent{apikey}
}

func (agent *OpenExchangeProxyAgent) GetRates() (RateResponse, error) {
	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", agent.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return RateResponse{}, err
	}

	defer resp.Body.Close()

	var incomingResponse IncomingRatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&incomingResponse); err != nil {
		return RateResponse{}, err
	}

	var outgoingResult RateResponse
	outgoingResult.Rates = incomingResponse.Rates
	outgoingResult.Base = incomingResponse.Base
	outgoingResult.Source = "openexchangerates.org"
	outgoingResult.Timestamp = time.Now()

	return outgoingResult, nil
}

// GetSymbols fetches and returns the currency symbols
func (agent *OpenExchangeProxyAgent) GetSymbols() (SymbolsResponse, error) {
	url := "https://openexchangerates.org/api/currencies.json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var incomingResponse IncomingSymbolsResponse
	if err := json.NewDecoder(resp.Body).Decode(&incomingResponse); err != nil {
		return nil, err
	}

	result := SymbolsResponse(incomingResponse)

	return result, nil
}
