package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/stanleyh24/restaurant-management/controllers"
)

func OrderItemRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Get("/orderItems", controller.GetOrderItems)
	incomingRoutes.Get("/orderItems/:orderItem_id", controller.GetOrderItem)
	incomingRoutes.Get("/orderItems/:order_id", controller.GetOrderItemsByOrder)
	incomingRoutes.Post("/orderItems/create", controller.CreateOrderItem)
	incomingRoutes.Patch("/orderItems/:orderItem_id", controller.UpdateOrderItem)
}
