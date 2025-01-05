package router

import (
	"github.com/YasenMakioui/gostore/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/gostore")

	api.Get("*", handler.ListObject)

	api.Post("*", handler.CreateObject)

	api.Delete("*", handler.DeleteOjbect)

}
