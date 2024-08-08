package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/secona/url-shortener/database"
)

func main() {
	log.Println("Running database migrations...")

	db := database.Open()

	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})

	if err != nil {
		log.Fatalf("Error creating sqlite3 instance: %s", err.Error())
	}

	migration, err := migrate.NewWithDatabaseInstance("file://migrations", "file:main.sqlite3", driver)

	if err != nil {
		log.Fatalf("Error creating migration instance: %s", err.Error())
	}

	if err := migration.Up(); err != nil {
		log.Fatalf("Error running up migration: %s", err.Error())
	}

	log.Println("Successfully run database migrations...")
}
