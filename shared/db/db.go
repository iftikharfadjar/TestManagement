package db

import (
	"log"

	"github.com/pocketbase/pocketbase"
)

var app *pocketbase.PocketBase

func Init() *pocketbase.PocketBase {
	app = pocketbase.New()
	return app
}

func GetApp() *pocketbase.PocketBase {
	if app == nil {
		log.Fatal("PocketBase is not initialized")
	}
	return app
}