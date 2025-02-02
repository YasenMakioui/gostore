package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type Disk struct {
	Size int
}

func GetDiskUsage(c *fiber.Ctx) error {
	return c.SendString("Disk usage")
}
