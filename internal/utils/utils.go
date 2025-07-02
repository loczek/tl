package utils

import "github.com/gofiber/fiber/v2"

func GetRequestId(c *fiber.Ctx) string {
	requestId := c.Locals("requestid").(string)
	return requestId
}
