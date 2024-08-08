package routing

import (
	"encoding/json"
	"net/http"
	"uacc-backend/contract"
)

func SetupLocalizations(
	mux *http.ServeMux,
) {
	mux.HandleFunc("/api/v1/localizations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		unique := make(map[string][]string)
		// Hungary
		unique["HUF"] = []string{"Ft"}
		// Bosnian
		unique["BAM"] = []string{"KM"}
		// Australia
		unique["AUD"] = []string{"AUD$"}
		// Brazil
		unique["BRL"] = []string{"R$"}
		// New Zealand
		unique["NZD"] = []string{"NZD$"}
		// USA
		unique["USD"] = []string{"USD$", "US$", "US $"}
		// Canada
		unique["CAD"] = []string{"CDN$"}
		// Mexico
		unique["MXN"] = []string{"MXN$"}
		// EU
		unique["EUR"] = []string{"€"}
		// UK
		unique["GBP"] = []string{"£", "￡"}
		// Japan
		unique["JPY"] = []string{"JP¥", "円"}
		// China
		unique["CNY"] = []string{"CN¥", "元"}
		// India
		unique["INR"] = []string{"₹", "Rs"}
		// Russia
		unique["RUB"] = []string{"₽"}
		// Kazakhstan
		unique["KZT"] = []string{"₸"}
		// Turkey
		unique["TRY"] = []string{"₺", "TL"}
		// Ukraine
		unique["UAH"] = []string{"₴"}
		// Thailand
		unique["THB"] = []string{"฿"}
		// Poland
		unique["PLN"] = []string{"zł"}
		// South Korea
		unique["KRW"] = []string{"₩"}
		// Bulgaria
		unique["BGN"] = []string{"лв"}
		// Czechia
		unique["CZK"] = []string{"Kč"}
		// South Africa
		//unique["ZAR"] = []string{ "R" } // this is just really annoying
		// Bitcoin
		unique["BTC"] = []string{"₿"}
		// Monero
		unique["XMR"] = []string{"ɱ"}
		// Ethereum
		unique["ETH"] = []string{"Ξ"}
		// Litecoin
		unique["LTC"] = []string{"Ł"}

		result := contract.LocalizationResponse{}
		result.Unique = unique

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})
}
