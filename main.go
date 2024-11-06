package main

import (
	"log"

	"github.com/aldiandyaIrsyad/uber-eats/config"
	"github.com/aldiandyaIrsyad/uber-eats/routes"
	"github.com/aldiandyaIrsyad/uber-eats/seeders"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to MongoDB
	client, ctx, cancel := config.ConnectDB()
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Seed database
	if err := seeders.SeedDatabase(client); err != nil {
		log.Printf("Error seeding database: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Setup routes with dependency injection
	routeHandler := routes.NewRouteHandler(client)
	routeHandler.SetupRoutes(r)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
