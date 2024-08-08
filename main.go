package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

type NotFoundData struct {
	Slug string
}

var slugRegex = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)

func main() {
	links := map[string]string{
		"google": "https://google.com",
		"pacil":  "https://cs.ui.ac.id",
	}

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

		http.Redirect(w, r, value, 301)
	})

	mux.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		slug := r.FormValue("slug")

		if !slugRegex.MatchString(slug) {
			fmt.Fprintf(w, "Shortened link must only contain alphabets, numbers, and hyphens!")
			return
		}

		url := r.FormValue("url")

		_, exists := links[slug]

		if exists {
			fmt.Fprintf(w, "Shortened link <strong>%s<strong> already exists!", slug)
			return
		}

		links[slug] = url
		fmt.Fprintf(w, "Successfully shortened link!")
	})

	http.ListenAndServe(":8080", mux)
}
