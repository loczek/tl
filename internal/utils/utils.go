package utils

import "github.com/gofiber/fiber/v2"

func GetRequestID(c *fiber.Ctx) string {
	requestID := c.Locals("requestid").(string)
	return requestID
}
