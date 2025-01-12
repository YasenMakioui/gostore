package main

import (
	"github.com/YasenMakioui/gostore/router"
	//"github.com/YasenMakioui/gostore/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	app := fiber.New()

	router.SetupRoutes(app)

	app.Use(cors.New())

	app.Listen(":3000")
}
