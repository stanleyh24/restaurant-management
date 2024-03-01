package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/stanleyh24/restaurant-management/controllers"
)

func OrderRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Get("/orders", controller.GetOrders)
	incomingRoutes.Get("/orders/:order_id", controller.GetOrder)
	incomingRoutes.Post("/orders/create", controller.CreateOrder)
	incomingRoutes.Patch("/orders/:order_id", controller.UpdateOrder)
}
