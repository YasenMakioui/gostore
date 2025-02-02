package main

import (
	"github.com/YasenMakioui/gostore/internal/middleware"
	"github.com/YasenMakioui/gostore/internal/router"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	router.SetupRoutes(app)

	middleware.SetupMiddlewares(app)

	app.Listen(":3000")
}
