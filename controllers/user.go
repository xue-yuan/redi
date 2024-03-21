package controllers

import (
	"fmt"
	"redi/constants"

	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	userID := c.Locals(constants.UserID)
	fmt.Println(userID)

	return constants.EmptyResponse(c)
}
