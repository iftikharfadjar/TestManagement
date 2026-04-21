package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"

	"boilerplate/graph"
)

func SetupRoutes(app *fiber.App) {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	app.Post("/query", adaptor.HTTPHandler(srv))

	playSrv := playground.Handler("GraphQL playground", "/query")
	app.Get("/", adaptor.HTTPHandler(playSrv))
}