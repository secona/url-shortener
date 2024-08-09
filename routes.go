package urlshortener

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/secona/url-shortener/database"
	"google.golang.org/api/idtoken"
)

func CreateMux(clientID string) *http.ServeMux {
	db := database.Open()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.String()[1:]

		if slug == "" {
			t := template.Must(template.ParseFiles("templates/index.html"))
			t.Execute(w, clientID)
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

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		credential := r.FormValue("credential")

		payload, err := idtoken.Validate(r.Context(), credential, clientID)

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		fmt.Fprintf(w, payload.Claims["email"].(string))
	})

	return mux
}
