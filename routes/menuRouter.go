package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/stanleyh24/restaurant-management/controllers"
)

func MenuRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Get("/menus", controller.GetMenus)
	incomingRoutes.Get("/menus/:menu_id", controller.GetMenu)
	incomingRoutes.Post("/menus/create", controller.CreateMenu)
	incomingRoutes.Patch("/menus/:menu_id", controller.UpdateMenu)
}
