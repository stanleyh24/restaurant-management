package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/stanleyh24/restaurant-management/database"
	"github.com/stanleyh24/restaurant-management/middleware"
	"github.com/stanleyh24/restaurant-management/routes"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	server := fiber.New()

	server.Use(cors.New())
	server.Use(logger.New())

	routes.UserRoutes(server)
	routes.Use(middleware.Authentication())
	routes.FoodRoutes(server)
	routes.MenuRoutes(server)
	routes.TableRoutes(server)
	routes.OrderRoutes(server)
	routes.OrderItemRoutes(server)
	routes.InvoiceRoutes(server)

	log.Fatal(server.Listen(":" + port))
}
