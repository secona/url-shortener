package urlshortener

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, err := r.Cookie("access_token")

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		token, err := jwt.ParseWithClaims(auth.Value, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})

		if err != nil {
			fmt.Fprintf(w, "Error verifying jwt: %s", err.Error())
			return
		}

		claims, ok := token.Claims.(*TokenClaims)

		if !ok {
			fmt.Fprintf(w, "Error getting token claims")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
