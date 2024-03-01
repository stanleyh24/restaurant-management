package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/stanleyh24/restaurant-management/controllers"
)

func FoodRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Get("/foods", controller.GetFoods)
	incomingRoutes.Get("/foods/:food_id", controller.GetFood)
	incomingRoutes.Post("/foods/create", controller.CreateFood)
	incomingRoutes.Patch("/foods/:food_id", controller.UpdateFood)
}
