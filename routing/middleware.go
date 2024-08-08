package routing

import (
	"net/http"
	"uacc-backend/util"
)

func SetupMiddleware(mux http.Handler) http.Handler {
	return validateApikey(mux)
}

func validateApikey(next http.Handler) http.Handler {
	apikey := util.GetOrCrash("ownApiKey")
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
