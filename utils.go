package urlshortener

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var slugRegex = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)

type TokenClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func parseURL(raw string) (*url.URL, error) {
	url, err := url.Parse(raw)

	if err != nil {
		return nil, err
	}

	// check if the scheme is either "http" or "https"
	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, fmt.Errorf("Invalid URL: must be either \"http\" or \"https\".")
	}

	// check if host is present
	if url.Host == "" {
		return nil, fmt.Errorf("Invalid URL: empty host.")
	}

	return url, nil
}

func parseSlug(slug string) (string, error) {
	if !slugRegex.MatchString(slug) {
		return "", fmt.Errorf("Shortened link must only contain alphabets, numbers, and hyphens!")
	}

	return slug, nil
}

func createAccessToken(userID int) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Hour * 24 * 7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})

	signed, err := token.SignedString(JwtSecret)

	return signed, expiresAt, err
}

func parseAccessToken(value string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(value, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)

	if !ok {
		return nil, fmt.Errorf("parsing token claims")
	}

	return claims, nil
}
