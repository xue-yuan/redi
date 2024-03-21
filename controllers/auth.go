package controllers

import (
	"context"
	"redi/config"
	"redi/constants"
	"redi/models"
	"redi/redis"
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
		return constants.UnauthorizedResponse(c, constants.ErrInvalidUsernameOrPassword)
	}

	t, err := utils.GenerateToken(user.UserID)
	if err != nil {
		return constants.ServerErrorResponse(c, err)
	}

	return constants.OkResponse(c, &fiber.Map{"token": t})
}

func Logout(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	userID := c.Locals(constants.UserID).(string)

	_, credentials := utils.GetAuthorizationSchemeAndParam(c.Get("Authorization"))
	redis.Client.Set(ctx, credentials, userID, config.Config.TokenTTL)

	return constants.EmptyResponse(c)
}

func RefreshToken(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	userID := c.Locals(constants.UserID).(string)

	_, credentials := utils.GetAuthorizationSchemeAndParam(c.Get("Authorization"))
	redis.Client.Set(ctx, credentials, userID, config.Config.OudatedTokenTTL)

	t, err := utils.GenerateToken(userID)
	if err != nil {
		return constants.ServerErrorResponse(c, err)
	}

	return constants.OkResponse(c, &fiber.Map{"token": t})
}
