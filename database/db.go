package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	DB *sql.DB
}

func Open() DB {
	db, err := sql.Open("sqlite3", "file:main.sqlite3")

	if err != nil {
		log.Fatalf("Error opening database: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %s", err.Error())
	}

	return DB{db}
}

func (db *DB) CreateShortenedLink(slug string, link string) error {
	_, err := db.DB.Exec("INSERT INTO links (slug, link) VALUES (?, ?)", slug, link)

	if err == sqlite3.ErrConstraint {
		return fmt.Errorf("Shortened link <strong>%s<strong> already exists!", slug)
	}

	return nil
}

func (db *DB) GetShortenedLink(slug string) (string, bool) {
	var link string
	if err := db.DB.QueryRow("SELECT link FROM links WHERE slug = ?", slug).Scan(&link); err != nil {
		return "", false
	}

	return link, true
}
