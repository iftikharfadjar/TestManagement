package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
)

func runMigrations(dbConn *sql.DB, dbType string, connString string) {
	log.Println("Running database migrations...")

	var driver database.Driver
	var err error

	if dbType == "postgres" {
		driver, err = postgres.WithInstance(dbConn, &postgres.Config{})
	} else if dbType == "sqlite" {
		driver, err = sqlite.WithInstance(dbConn, &sqlite.Config{})
	}

	if err != nil {
		log.Fatalf("could not instantiate migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./sql/migrations", dbType, driver)
	if err != nil {
		log.Fatalf("migration init failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("an error occurred while syncing the database.. %v", err)
	}
	log.Println("Database migrated successfully!")
}
