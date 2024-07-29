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

type InvoiceViewFormat struct {
	Invoice_id       string
	Payment_method   string
	Order_id         string
	Payment_status   *string
	Payment_due      any
	Table_number     any
	Payment_due_date time.Time
	Order_details    any
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	result, err := invoiceCollection.Find(context.TODO(), bson.M{})
	defer cancel()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error occured while listing invoices items"})
	}
	var allInvoices []bson.M

	if err := result.All(ctx, &allInvoices); err != nil {
		log.Fatal(err)
	}
	return c.Status(http.StatusOK).JSON(allInvoices)
}

func GetInvoice(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	invoiceId := c.Params("invoice_id")
	var invoice models.Invoice

	err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
	defer cancel()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error occured while listing invoice item"})
	}

	var invoiceView InvoiceViewFormat

	allOrderItems, err := ItemsByOrder(invoice.Order_id)
	invoiceView.Order_id = invoice.Order_id
	invoiceView.Payment_due_date = invoice.Payment_due_date

	invoiceView.Payment_method = "null"
	if invoice.Payment_method == nil {
		invoiceView.Payment_method = *invoice.Payment_method
	}

	invoiceView.Invoice_id = invoice.Invoice_id
	invoiceView.Payment_status = *&invoice.Payment_status
	invoiceView.Payment_due = allOrderItems[0]["payment_due"]
	invoiceView.Table_number = allOrderItems[0]["table_number"]
	invoiceView.Order_details = allOrderItems[0]["order_items"]

	return c.Status(http.StatusOK).JSON(invoiceView)
}

func CreateInvoice(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var invoice models.Invoice

	if err := c.BodyParser(&invoice); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var order models.Order

	err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
	defer cancel()

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "order was not found"})
	}

	status := "PENDING"
	if invoice.Payment_status == nil {
		invoice.Payment_status = &status
	}
	invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.ID = primitive.NewObjectID()
	invoice.Invoice_id = invoice.ID.Hex()

	validationErr := validate.Struct(invoice)

	if validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Error()})
	}

	result, insertErr := invoiceCollection.InsertOne(ctx, invoice)
	if insertErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "food item was not created"})
	}
	defer cancel()
	return c.Status(http.StatusOK).JSON(result)

}

func UpdateInvoice(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var invoice models.Invoice
	invoiceID := c.Params("invoice_id")

	if err := c.BodyParser(&invoice); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var updateObj primitive.D

	if invoice.Payment_method != nil {
		updateObj = append(updateObj, bson.E{"payment_method", invoice.Payment_method})
	}

	if invoice.Payment_status != nil {
		updateObj = append(updateObj, bson.E{"payment_status", invoice.Payment_status})
	}

	invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", invoice.Updated_at})

	upsert := true
	filter := bson.M{"invoice_id": invoiceID}

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	status := "PENDING"

	if invoice.Payment_status == nil {
		invoice.Payment_status = &status
	}
	result, err := invoiceCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "invoice item update failed"})
	}

	defer cancel()
	return c.Status(http.StatusOK).JSON(result)
}
