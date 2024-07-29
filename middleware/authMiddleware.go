package middleware

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/stanleyh24/restaurant-management/helpers"
)

func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/users/signup" || c.Path() == "/users/login" {
			return c.Next()
		}
		clientToken := c.GetRespHeader("token")
		fmt.Println(clientToken)
		if clientToken == "" {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "no Authorization header provided"})
		}
		claims, err := helpers.ValidateToken(clientToken)

		if err != "" {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)
		return c.Next()
	}
}
