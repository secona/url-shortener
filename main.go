package main

import (
	"html/template"
	"net/http"
)

type NotFoundData struct {
	Slug string
}

func main() {
	links := map[string]string{
		"google": "https://google.com",
		"pacil":  "https://cs.ui.ac.id",
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.String()[1:]
		value, ok := links[slug]

		if !ok {
			t := template.Must(template.ParseFiles("templates/404.html"))
			t.Execute(w, slug)
			return
		}

		http.Redirect(w, r, value, 301)
	})

	http.ListenAndServe(":8080", mux)
}
