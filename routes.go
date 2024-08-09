package urlshortener

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/secona/url-shortener/database"
	"google.golang.org/api/idtoken"
)

func CreateMux(clientID string, jwtSecret string) *chi.Mux {
	db := database.Open()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("templates/index.html"))
		t.Execute(w, clientID)
	})

	r.Post("/shorten", func(w http.ResponseWriter, r *http.Request) {
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

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		credential := r.FormValue("credential")

		payload, err := idtoken.Validate(r.Context(), credential, clientID)

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		userID, err := db.UpsertUser(database.User{
			Name:  payload.Claims["name"].(string),
			Email: payload.Claims["email"].(string),
			Pic:   payload.Claims["picture"].(string),
		})

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		exp := time.Now().Add(7 * 24 * time.Hour)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
			"exp":     exp,
		})
		signed, err := token.SignedString([]byte(jwtSecret))

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    signed,
			Expires:  exp,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		fmt.Fprintf(w, payload.Claims["email"].(string))
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.String()[1:]
		value, ok := db.GetShortenedLink(slug)

		if !ok {
			t := template.Must(template.ParseFiles("templates/404.html"))
			t.Execute(w, slug)
			return
		}

		http.Redirect(w, r, value, 301)
	})

	return r
}
