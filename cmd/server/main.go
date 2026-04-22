package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"boilerplate/services/auth/delivery/graphql"
	authRest "boilerplate/services/auth/delivery/rest"
	testRest "boilerplate/services/test/delivery/rest"
	"boilerplate/services/auth/usecase"
	testUsecase "boilerplate/services/test/usecase"
	"boilerplate/shared/config"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Booting with DB_TYPE: %s", cfg.DBType)

	repos := DbConnSwitcher(cfg)

	authUseCase := usecase.NewAuthUseCase(repos.Auth)
	testUseCase := testUsecase.NewTestUseCase(repos.Test)

	app := fiber.New()

	authHandler := authRest.NewAuthHandler(authUseCase)
	authHandler.SetupRoutes(app)

	testHandler := testRest.NewTestHandler(testUseCase)
	testHandler.SetupRoutes(app)

	graphql.SetupRoutes(app)

	if cfg.DBType == "pocketbase" {
		app.All("/_/*", proxy.Forward("http://127.0.0.1:8090/_/"))
		app.All("/api/*", proxy.Forward("http://localhost:8090/api/"))
	}

	log.Printf("Starting Fiber Server on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Fiber server failed: %v", err)
	}
}
