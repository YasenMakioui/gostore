package middleware

import "github.com/gofiber/fiber/v2"

func SetupMiddlewares(app *fiber.App) {
	// Call middleware functions
	Cors(app)
	Logger(app)
}
