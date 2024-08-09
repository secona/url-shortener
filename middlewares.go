package urlshortener

import (
	"context"
	"fmt"
	"net/http"
)

func authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, err := r.Cookie("access_token")

		if err != nil {
			fmt.Fprintf(w, "Error getting access token: %s", err.Error())
			return
		}

		parsed, err := parseAccessToken(auth.Value)

		if err != nil {
			fmt.Fprintf(w, "Error verifying JWT: %s", err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", parsed.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
