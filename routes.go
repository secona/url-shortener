package urlshortener

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/secona/url-shortener/database"
	"google.golang.org/api/idtoken"
)

var ClientID string
var JwtSecret []byte

func CreateMux() *chi.Mux {
	database.Open()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("templates/index.html"))
		t.Execute(w, ClientID)
	})

	r.With(authenticated).Post("/shorten", func(w http.ResponseWriter, r *http.Request) {
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

		_, err = database.Link{
			Slug: slug,
			Link: url.String(),
			UserID: r.Context().Value("user_id").(int),
		}.Create()

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		fmt.Fprintf(w, "Successfully shortened link!")
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		credential := r.FormValue("credential")

		payload, err := idtoken.Validate(r.Context(), credential, ClientID)

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		user, err := database.User{
			Name:  payload.Claims["name"].(string),
			Email: payload.Claims["email"].(string),
			Pic:   payload.Claims["picture"].(string),
		}.Upsert()

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		token, expiresAt, err := createAccessToken(user.ID)

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    token,
			Expires:  expiresAt,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		fmt.Fprintf(w, payload.Claims["email"].(string))
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.String()[1:]
		link, ok := database.Link{Slug: slug}.Get()

		if !ok {
			t := template.Must(template.ParseFiles("templates/404.html"))
			t.Execute(w, slug)
			return
		}

		http.Redirect(w, r, link.Link, 301)
	})

	return r
}
