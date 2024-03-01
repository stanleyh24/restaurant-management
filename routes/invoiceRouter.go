package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/stanleyh24/restaurant-management/controllers"
)

func InvoiceRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Get("/invoices", controller.GetInvoices)
	incomingRoutes.Get("/invoices/:invoice_id", controller.GetInvoice)
	incomingRoutes.Post("/invoices/create", controller.CreateInvoice)
	incomingRoutes.Patch("/invoices/:invoice_id", controller.UpdateInvoice)
}
