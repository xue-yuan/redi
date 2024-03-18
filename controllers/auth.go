package controllers

import (
	"context"
	"errors"
	"redi/constants"
	"redi/models"
	"redi/utils"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	user := &models.User{}

	if err := c.BodyParser(user); err != nil {
		return constants.ClientErrorResponse(c, err)
	}

	if err := user.Create(ctx); err != nil {
		return constants.ServerErrorResponse(c, err)
	}

	t, err := utils.GenerateToken(user.UserID)
	if err != nil {
		return constants.ServerErrorResponse(c, err)
	}

	return constants.OkResponse(c, &fiber.Map{"token": t})
}

func Login(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	user := &models.User{}

	if err := c.BodyParser(user); err != nil {
		return constants.ClientErrorResponse(c, err)
	}

	ok, err := user.Login(ctx, user.Password)
	if err != nil {
		return constants.ServerErrorResponse(c, err)
	} else if !ok {
		return constants.UnauthorizedResponse(c, errors.New("invalid username or password"))
	}

	t, err := utils.GenerateToken(user.UserID)
	if err != nil {
		return constants.ServerErrorResponse(c, err)
	}

	return constants.OkResponse(c, &fiber.Map{"token": t})
}

func Logout(c *fiber.Ctx) error {
	return nil
}

func RefreshToken(c *fiber.Ctx) error {
	return nil
}
