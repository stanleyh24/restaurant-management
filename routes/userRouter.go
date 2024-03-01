package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/stanleyh24/restaurant-management/controllers"
)

func UserRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Get("/users", controller.GetUsers)
	incomingRoutes.Get("/users/:user_id", controller.GetUser)
	incomingRoutes.Post("/users/signup", controller.SignUp)
	incomingRoutes.Get("/users/login", controller.Login)
}
