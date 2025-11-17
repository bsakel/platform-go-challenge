package main

import (
	"context"
	"platform-go-challenge/api"
	"platform-go-challenge/db"
	"platform-go-challenge/graph"
	"platform-go-challenge/graph/resolvers"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func graphqlHandler(resolver graph.ResolverRoot) gin.HandlerFunc {
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
	resolver := &resolvers.Resolver{
		DB: db.GormDB,
	}

	//REST routes
	// Audience routes
	router.POST("/audience", api.CreateAudience)
	router.GET("/audiences", api.GetAudiences)
	router.GET("/audience/:id", api.GetAudience)
	router.PUT("/audience/:id", api.UpdateAudience)
	router.DELETE("/audience/:id", api.DeleteAudience)

	// Chart routes
	router.POST("/chart", api.CreateChart)
	router.GET("/charts", api.GetCharts)
	router.GET("/chart/:id", api.GetChart)
	router.PUT("/chart/:id", api.UpdateChart)
	router.DELETE("/chart/:id", api.DeleteChart)

	// Insight routes
	router.POST("/insight", api.CreateInsight)
	router.GET("/insights", api.GetInsights)
	router.GET("/insight/:id", api.GetInsight)
	router.PUT("/insight/:id", api.UpdateInsight)
	router.DELETE("/insight/:id", api.DeleteInsight)

	// UserFavourite routes
	router.POST("/userfavourite", api.CreateUserFavourite)
	router.GET("/userfavourites", api.GetUserFavourites)
	router.GET("/userfavourite/:id", api.GetUserFavourite)
	router.GET("/userfavourites/user/:userid", api.GetUserFavouritesByUser)
	router.PUT("/userfavourite/:id", api.UpdateUserFavourite)
	router.DELETE("/userfavourite/:id", api.DeleteUserFavourite)

	// GraphQL routes
	router.POST("/graphql", graphqlHandler(resolver))
	router.GET("/graphql", playgroundHandler())

	router.Run(":8080")
}
