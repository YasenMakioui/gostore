package router

import (
	. "github.com/YasenMakioui/gostore/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1/gostore")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Root")
	})

	api.Get("/store/*", GetObject)
	api.Post("/store/*", CreateObject)
	// api.Delete("/store/*", DeleteOjbect)
	// api.Put("/store/*", ModifyObject)

	api.Get("/stats/usage/disk", GetDiskUsage)

}
