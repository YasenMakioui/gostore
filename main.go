package main

import (
	"github.com/YasenMakioui/gostore/middleware"
	"github.com/YasenMakioui/gostore/router"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	router.SetupRoutes(app)

	middleware.SetupMiddlewares(app)

	app.Listen(":3000")
}
