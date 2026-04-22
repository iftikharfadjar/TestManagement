package main

import (
	"database/sql"
	"log"

	authDomain "boilerplate/services/auth/domain"
	testDomain "boilerplate/services/test/domain"
	"boilerplate/shared/adapter/pocketbase"
	sqliteAdapter "boilerplate/shared/adapter/sqlite_adapter"
	"boilerplate/shared/config"
	"boilerplate/shared/db"
)

type Repositories struct {
	Auth authDomain.AuthRepository
	Test testDomain.TestRepository
}

func DbConnSwitcher(cfg *config.Config) *Repositories {
	repos := &Repositories{}

	switch cfg.DBType {
	case "postgres":
		dbConn, err := sql.Open("postgres", cfg.DBConnString)
		if err != nil {
			log.Fatal(err)
		}
		runMigrations(dbConn, "postgres", cfg.DBConnString)
		repos.Auth = sqliteAdapter.NewAuthRepository(dbConn, cfg.JWTSecret)
		repos.Test = sqliteAdapter.NewTestRepository(dbConn)

	case "sqlite":
		dbConn, err := sql.Open("sqlite", cfg.DBConnString)
		if err != nil {
			log.Fatal(err)
		}
		runMigrations(dbConn, "sqlite", cfg.DBConnString)
		repos.Auth = sqliteAdapter.NewAuthRepository(dbConn, cfg.JWTSecret)
		repos.Test = sqliteAdapter.NewTestRepository(dbConn)

	case "pocketbase":
		fallthrough
	default:
		pbApp := db.Init()
		repos.Auth = pocketbase.NewAuthRepository(pbApp)

		go func() {
			if err := pbApp.Start(); err != nil {
				log.Fatalf("PocketBase start failed: %v", err)
			}
		}()
	}
	return repos
}
