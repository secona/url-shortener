package urlshortener

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/secona/url-shortener/database"
)

func CreateMux() *http.ServeMux {
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
		slug, err := parseSlug(r.FormValue("slug"))

		if err != nil {
			fmt.Fprintf(w, err.Error())
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

	return mux
}
