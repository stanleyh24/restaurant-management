package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stanleyh24/restaurant-management/database"
	"github.com/stanleyh24/restaurant-management/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func GetFoods(c *fiber.Ctx) error {
	return nil
}

func GetFood(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	foodId := c.Params("food_id")
	var food models.Food

	err := foodCollection.FindOne(ctx, bson.M{"food": foodId}).Decode(&food)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error occured while fetching food item"})
	}
	return c.Status(http.StatusOK).JSON(food)
}

func CreateFood(c *fiber.Ctx) error {
	return nil
}

func UpdateFood(c *fiber.Ctx) error {
	return nil
}

func round(num float64) int {
	return 0
}

func toFixed(num float64, precision int) float64 {
	return 0.1
}
