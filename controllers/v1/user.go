package v1

import (
	"redi/constants"

	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	return constants.EmptyResponse(c)
}
