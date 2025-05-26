package jwt

import (
	"context"
	"net/http"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenStr, err := extractToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := validateToken(tokenStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := extractUserID(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
