package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stanleyh24/restaurant-management/database"
	"github.com/stanleyh24/restaurant-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	result, err := orderCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error occured while getting orders items"})
	}

	var allOrders []bson.M

	if err = result.All(ctx, &allOrders); err != nil {
		log.Fatal(err)
	}

	return c.Status(http.StatusOK).JSON(allOrders)
}

func GetOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	orderId := c.Params("order_id")
	var order models.Order

	err := orderCollection.FindOne(ctx, bson.M{"order": orderId}).Decode(&order)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error occured while fetching order item"})
	}
	return c.Status(http.StatusOK).JSON(order)
}

func CreateOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var table models.Table
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	validationErr := validate.Struct(order)

	if validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Error()})
	}

	if order.Table_id != nil {
		err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
		defer cancel()

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "table was not found"})
		}
	}

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	result, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "order item was not created"})
	}
	return c.Status(http.StatusOK).JSON(result)
}

func UpdateOrder(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var table models.Table
	var order models.Order
	var updateObj primitive.D
	orderID := c.Params("order_id")
	if err := c.BodyParser(&order); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if order.Table_id != nil {
		err := orderCollection.FindOne(ctx, bson.M{"table_id": orderID}).Decode(&table)
		defer cancel()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Menu was not found"})
		}
		updateObj = append(updateObj, bson.E{"menu", order.Table_id})
	}

	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", order.Updated_at})

	upsert := true
	filter := bson.M{"order_id": orderID}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	result, err := orderCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "order update failed"})
	}
	defer cancel()

	return c.Status(http.StatusOK).JSON(result)
}

func OrderItemOrderCreator(order models.Order) string {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)
	defer cancel()

	return order.Order_id
}
