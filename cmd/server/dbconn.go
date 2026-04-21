package main

import (
	"database/sql"
	"log"

	"boilerplate/services/auth/domain"
	"boilerplate/shared/adapter/pocketbase"
	sqliteAdapter "boilerplate/shared/adapter/sqlite_adapter"
	"boilerplate/shared/config"
	"boilerplate/shared/db"
)

func DbConnSwitcher(cfg *config.Config) domain.AuthRepository {
	var authRepo domain.AuthRepository
	switch cfg.DBType {
	case "postgres":
		dbConn, err := sql.Open("postgres", cfg.DBConnString)
		if err != nil {
			log.Fatal(err)
		}
		runMigrations(dbConn, "postgres", cfg.DBConnString)
		authRepo = sqliteAdapter.NewAuthRepository(dbConn, cfg.JWTSecret)

	case "sqlite":
		dbConn, err := sql.Open("sqlite", cfg.DBConnString)
		if err != nil {
			log.Fatal(err)
		}
		runMigrations(dbConn, "sqlite", cfg.DBConnString)
		authRepo = sqliteAdapter.NewAuthRepository(dbConn, cfg.JWTSecret)

	case "pocketbase":
		fallthrough
	default:
		pbApp := db.Init()
		authRepo = pocketbase.NewAuthRepository(pbApp)

		go func() {
			if err := pbApp.Start(); err != nil {
				log.Fatalf("PocketBase start failed: %v", err)
			}
		}()
	}
	return authRepo
}
