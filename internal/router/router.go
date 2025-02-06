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

	api.Get("/store/*", GetFilesystemObject)
	api.Post("/store/*", CreateFilesystemObject)
	api.Delete("/store/*", DeleteFilesystemOjbect)
	api.Put("/store/*", ModifyFilesystemObject)

	//api.Get("/stats", GetStats)
	//api.Get("/stats/backups", GetBackupStats) // will return the backup stats like running backups and so on

	api.Get("/backup/*", GetBackups)
	//api.Post("/backup/*", PostBackup)
	//api.Delete("/backup/*", DeleteBackup)
}
