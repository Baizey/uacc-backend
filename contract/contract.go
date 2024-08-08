package contract

type LocalizationResponse struct {
	Unique map[string][]string `json:"unique"`
}

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

type RatesResponse struct {
	Rates []RateResponse `json:"rates"`
}

type Data struct {
	Symbols map[string]string
	// map[from][to]rate
	Rates map[string]map[string]RateResponse
	// map[to][]rates
	NewRates map[string]RatesResponse
}
