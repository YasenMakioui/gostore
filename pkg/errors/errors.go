package errors

import "github.com/gofiber/fiber/v2"

func FormatError(msg string) fiber.Map {
	// used to give a format returned to a userç
	// DRY pragmatic -> If we want to change the err format we just have to modify this helper
	return fiber.Map{
		"error": msg,
	}
}
