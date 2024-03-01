package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrderItems(c *fiber.Ctx) error {
	return nil
}

func GetOrderItem(c *fiber.Ctx) error {
	return nil
}

func GetOrderItemsByOrder(c *fiber.Ctx) error {
	return nil
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {

}

func CreateOrderItem(c *fiber.Ctx) error {
	return nil
}

func UpdateOrderItem(c *fiber.Ctx) error {
	return nil
}
