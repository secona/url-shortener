package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"regexp"

	"github.com/secona/url-shortener/database"
)

var slugRegex = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)

func main() {
	db := database.Open()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.String()[1:]

		if slug == "" {
			t := template.Must(template.ParseFiles("templates/index.html"))
			t.Execute(w, nil)
			return
		}

		value, ok := db.GetShortenedLink(slug)

		if !ok {
			t := template.Must(template.ParseFiles("templates/404.html"))
			t.Execute(w, slug)
			return
		}

		http.Redirect(w, r, value, 301)
	})

	mux.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		slug := r.FormValue("slug")

		if !slugRegex.MatchString(slug) {
			fmt.Fprintf(w, "Shortened link must only contain alphabets, numbers, and hyphens!")
			return
		}

		url, err := parseURL(r.FormValue("url"))

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		err = db.CreateShortenedLink(slug, url.String())

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		fmt.Fprintf(w, "Successfully shortened link!")
	})

	http.ListenAndServe(":8080", mux)
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
