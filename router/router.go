package router

import (
	"github.com/YasenMakioui/gostore/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1/gostore")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Root")
	})

	api.Get("/store/*", handler.GetObject)
	api.Post("/store/*", handler.CreateObject)
	api.Delete("/store/*", handler.DeleteOjbect)
	api.Put("/store/*", handler.ModifyObject)

	api.Get("/stats/usage/disk", handler.GetDiskUsage)

}
