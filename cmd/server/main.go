package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/secona/url-shortener"
)

func init() {
	env, err := godotenv.Read()

	if err != nil {
		log.Fatalf("Error reading .env file: %s", err.Error())
	}

	urlshortener.ClientID = env["GOOGLE_CLIENT_ID"]
	urlshortener.JwtSecret = []byte(env["JWT_SECRET"])
}

func main() {
	mux := urlshortener.CreateMux()
	http.ListenAndServe(":8080", mux)
}
