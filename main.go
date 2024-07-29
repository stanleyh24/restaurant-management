package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	//customMidleware "github.com/stanleyh24/restaurant-management/middleware"
	"github.com/stanleyh24/restaurant-management/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	server := fiber.New()

	server.Use(cors.New())
	server.Use(logger.New())
	routes.UserRoutes(server)
	//server.Use(customMidleware.Authentication())
	routes.FoodRoutes(server)
	routes.MenuRoutes(server)
	routes.TableRoutes(server)
	routes.OrderRoutes(server)
	routes.OrderItemRoutes(server)
	routes.InvoiceRoutes(server)

	log.Fatal(server.Listen(":" + port))
}
