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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	result, err := orderCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error occured while getting orders items"})
	}

	var allTables []bson.M

	if err = result.All(ctx, &allTables); err != nil {
		log.Fatal(err)
	}

	return c.Status(http.StatusOK).JSON(allTables)
}

func GetTable(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	tableId := c.Params("table_id")
	var table models.Table

	err := tableCollection.FindOne(ctx, bson.M{"table": tableId}).Decode(&table)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error occured while fetching the table "})
	}
	return c.Status(http.StatusOK).JSON(table)
}

func CreateTable(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var table models.Table

	if err := c.BodyParser(&table); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	validationErr := validate.Struct(table)

	if validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Error()})
	}

	table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	table.ID = primitive.NewObjectID()
	table.Table_id = table.ID.Hex()

	result, err := tableCollection.InsertOne(ctx, table)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "table  was not created"})
	}
	defer cancel()
	return c.Status(http.StatusOK).JSON(result)
}

func UpdateTable(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var table models.Table
	tableID := c.Params("table_id")
	if err := c.BodyParser(&table); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var updateObj primitive.D

	if table.Number_of_guests != nil {
		updateObj = append(updateObj, bson.E{"number_of_guests", table.Number_of_guests})
	}

	if table.Table_number != nil {
		updateObj = append(updateObj, bson.E{"table_number", table.Table_number})
	}

	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	upsert := true
	filter := bson.M{"table_id": tableID}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := tableCollection.UpdateOne(
		ctx,
		filter,
		bson.D{{"$set", updateObj}},
		&opt,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "table item update failed"})
	}
	defer cancel()

	return c.Status(http.StatusOK).JSON(result)
}
