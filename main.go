package main

import (
	"context"
	"platform-go-challenge/api"
	"platform-go-challenge/db"
	"platform-go-challenge/graph"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func graphqlHandler(resolver *graph.Resolver) gin.HandlerFunc {
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: resolver})

	h := handler.New(schema)

	// Add DataLoader middleware
	h.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		ctx = graph.DataLoaderMiddleware(db.GormDB)(ctx)
		return next(ctx)
	})

	// Enable Automatic Persisted Queries (APQ) for caching
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100), // Cache up to 100 queries
	})

	// Add introspection caching
	h.Use(extension.Introspection{})

	// Configure transports
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.GET{})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	db.InitDB()
	router := gin.Default()

	// GraphQL resolver
	resolver := &graph.Resolver{
		DB: db.GormDB,
	}

	//REST routes
	router.POST("/audience", api.CreateAudience)
	router.GET("/audiences", api.GetAudiences)
	router.GET("/audience/:id", api.GetAudience)
	router.PUT("/audience/:id", api.UpdateAudience)
	router.DELETE("/audience/:id", api.DeleteAudience)

	// GraphQL routes
	router.POST("/graphql", graphqlHandler(resolver))
	router.GET("/graphql", playgroundHandler())

	router.Run(":8080")
}
