package constants

import (
	"github.com/gofiber/fiber/v2"
)

func ServerErrorResponse(c *fiber.Ctx, msg error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": msg.Error()})
}

func ClientErrorResponse(c *fiber.Ctx, msg error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg.Error()})
}

func BadRequestResponse(c *fiber.Ctx, msg error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg.Error()})
}

func UnauthorizedResponse(c *fiber.Ctx, msg error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": msg.Error()})
}

func ForbiddenResponse(c *fiber.Ctx, msg error) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": msg.Error()})
}

func OkResponse(c *fiber.Ctx, data *fiber.Map) error {
	return c.Status(fiber.StatusOK).JSON(data)
}
