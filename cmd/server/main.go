package main

import (
	"net/http"

	"github.com/secona/url-shortener"
)

func main() {
	mux := urlshortener.CreateMux()
	http.ListenAndServe(":8080", mux)
}
