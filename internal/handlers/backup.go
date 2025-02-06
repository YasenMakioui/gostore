package handlers

import (
	. "github.com/YasenMakioui/gostore/internal/service"
	"github.com/gofiber/fiber/v2"
)

func GetBackups(requestContext *fiber.Ctx) error {

	backup := NewBackup("test", "test", "test")

	return requestContext.JSON(backup)
}
