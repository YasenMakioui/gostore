package main

import (
	"github.com/YasenMakioui/gostore/internal/middleware"
	"github.com/YasenMakioui/gostore/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()

	middleware.Logger(app)
	middleware.Cors(app)

	router.SetupRoutes(app)

	app.Use(logger.New())

	app.Listen(":3000")
}
