package handlers

import "github.com/gofiber/fiber/v2"

func GetStats(requestContext *fiber.Ctx) error {
	return requestContext.SendString("stats")
}
