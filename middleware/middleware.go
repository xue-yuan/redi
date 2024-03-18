package middleware

import (
	"context"
	"redi/constants"
	"redi/database"

	"github.com/gofiber/fiber/v2"
)

func SetupContext(c *fiber.Ctx) error {
	ctx := context.WithValue(c.Context(), constants.DB, database.Pool)
	c.Locals(constants.CTX, ctx)
	return c.Next()
}
