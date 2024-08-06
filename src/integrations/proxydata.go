package integrations

import "time"

type RateResponse struct {
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
	Timestamp time.Time          `json:"timestamp"`
	Source    string             `json:"source"`
}

// SymbolsResponse represents a map of symbols.
type SymbolsResponse map[string]string

// ProxyAgent defines the methods for fetching rates and symbols.
type ProxyAgent interface {
	GetRates() (RateResponse, error)
	GetSymbols() (SymbolsResponse, error)
}
