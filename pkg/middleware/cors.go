package middleware

import (
	"net/http"
	"slices"
)

func CorsMiddleware(next http.Handler, corsAllowList []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if slices.Contains(corsAllowList, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Block if origin or IP is not in the allow list
		if !slices.Contains(corsAllowList, origin) {
			http.Error(w, "Forbidden: Access denied", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
