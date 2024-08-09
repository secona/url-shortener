package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/secona/url-shortener"
)

func main() {
	env, err := godotenv.Read()

	if err != nil {
		log.Fatalln("Error reading .env!")
	}

	mux := urlshortener.CreateMux(env["GOOGLE_CLIENT_ID"], env["JWT_SECRET"])
	http.ListenAndServe(":8080", mux)
}
