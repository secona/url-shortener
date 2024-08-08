package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
)

type NotFoundData struct {
	Slug string
}

var slugRegex = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)

func main() {
	u, _ := url.Parse("https://google.com")
	fmt.Printf("%+v\n", *u)

	links := map[string]url.URL{}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.String()[1:]

		if slug == "" {
			t := template.Must(template.ParseFiles("templates/index.html"))
			t.Execute(w, nil)
			return
		}

		value, ok := links[slug]

		if !ok {
			t := template.Must(template.ParseFiles("templates/404.html"))
			t.Execute(w, slug)
			return
		}

		http.Redirect(w, r, value.String(), 301)
	})

	mux.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		slug := r.FormValue("slug")

		if !slugRegex.MatchString(slug) {
			fmt.Fprintf(w, "Shortened link must only contain alphabets, numbers, and hyphens!")
			return
		}

		_, exists := links[slug]

		if exists {
			fmt.Fprintf(w, "Shortened link <strong>%s<strong> already exists!", slug)
			return
		}

		url, err := parseURL(r.FormValue("url"))

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		links[slug] = *url

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
