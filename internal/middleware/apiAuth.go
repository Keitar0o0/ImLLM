package middleware

import (
	"context"
	"github.com/Xarth-Mai/ImLLM/internal/utils"
	"net/http"
	"strings"
)

// ApiAuth is a middleware that checks the apikey and wraps the request context with the username
func ApiAuth(next http.HandlerFunc, userPasswd map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")

		//apiKey == Username + ";" + Passwd
		split := strings.Split(apiKey, ";")
		if len(split) != 2 {
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}
		username := split[0]

		if utils.ValidateMap(userPasswd, username, split[1]) {
			ctx := context.WithValue(r.Context(), "username", username)
			next(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}
