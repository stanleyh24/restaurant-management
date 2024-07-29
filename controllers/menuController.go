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

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	result, err := menuCollection.Find(context.TODO(), bson.M{})

	defer cancel()

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error occurred while listing menus items"})
	}
	var allMenus []bson.M

	if err = result.All(ctx, &allMenus); err != nil {
		log.Fatal(err)
	}
	return c.Status(http.StatusOK).JSON(allMenus)
}

func GetMenu(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	menuId := c.Params("menu_id")
	var menu models.Menu

	err := menuCollection.FindOne(ctx, bson.M{"food": menuId}).Decode(&menu)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error occured while fetching the menu "})
	}
	return c.Status(http.StatusOK).JSON(menu)
}

func CreateMenu(c *fiber.Ctx) error {
	var menu models.Menu
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	if err := c.BodyParser(&menu); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	validationErr := validate.Struct(menu)

	if validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Error()})
	}

	menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.ID = primitive.NewObjectID()
	menu.Menu_id = menu.ID.Hex()

	result, insertErr := menuCollection.InsertOne(ctx, menu)

	if insertErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "menu item was not created"})
	}
	defer cancel()
	return c.Status(http.StatusOK).JSON(result)
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(check) && end.After(start)
}

func UpdateMenu(c *fiber.Ctx) error {
	var menu models.Menu
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	if err := c.BodyParser(&menu); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	validationErr := validate.Struct(menu)

	if validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Error()})
	}

	menuId := c.Params("menu_id")
	filter := bson.M{"menu_id": menuId}
	var updateObj primitive.D

	if menu.Start_date != nil && menu.End_date != nil {
		if !inTimeSpan(*menu.Start_date, *menu.End_date, time.Now()) {
			defer cancel()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "kindly retype the time"})
		}

		updateObj = append(updateObj, bson.E{"start_date", menu.Start_date})
		updateObj = append(updateObj, bson.E{"end_date", menu.End_date})

		if menu.Name != "" {
			updateObj = append(updateObj, bson.E{"name", menu.Name})
		}
		if menu.Category != "" {
			updateObj = append(updateObj, bson.E{"category", menu.Category})
		}
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", menu.Updated_at})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := menuCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "menu update failed"})
		}

		defer cancel()
		return c.Status(http.StatusOK).JSON(result)
	}
	defer cancel()
	return nil

}
