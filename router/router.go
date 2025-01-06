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

	// add option to create a file or dir using querystring
	api.Post("*", handler.CreateObject)

	api.Delete("*", handler.DeleteOjbect)

}
