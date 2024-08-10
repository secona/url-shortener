package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Open() {
	var err error
	DB, err = sql.Open("sqlite3", "file:main.sqlite3")

	if err != nil {
		log.Fatalf("Error opening database: %s", err.Error())
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %s", err.Error())
	}
}

type Link struct {
	ID     string
	Slug   string
	Link   string
	UserID int
}

func (l Link) Create() (*Link, error) {
	_, err := DB.Exec(
		"INSERT INTO links (slug, link, user_id) VALUES (?, ?, ?)",
		l.Slug,
		l.Link,
		sql.NullInt32{
			Int32: int32(l.UserID),
			Valid: l.UserID != 0,
		},
	)

	if err != nil {
		if err, ok := err.(sqlite3.Error); ok {
			if err.Code == sqlite3.ErrConstraint {
				return nil, fmt.Errorf("Shortened link <strong>%s<strong> already exists!", l.Slug)
			}
		}

		return nil, fmt.Errorf("Error shortening link")
	}

	return &l, nil
}

func (l Link) Get() (*Link, bool) {
	if err := DB.QueryRow("SELECT link FROM links WHERE slug = ?", l.Slug).Scan(&l.Link); err != nil {
		return nil, false
	}

	return &l, true
}

type User struct {
	ID    int
	Name  string
	Email string
	Pic   string
}

func (u User) Upsert() (*User, error) {
	row := DB.QueryRow(
		"INSERT INTO users (name, email, pic) VALUES (?, ?, ?) ON CONFLICT DO UPDATE SET pic = excluded.pic RETURNING id",
		u.Name,
		u.Email,
		u.Pic,
	)

	if err := row.Scan(&u.ID); err != nil {
		return nil, err
	}

	return &u, nil
}
