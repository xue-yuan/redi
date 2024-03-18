package controllers

import (
	"context"
	"redi/constants"
	"redi/models"

	"github.com/gofiber/fiber/v2"
)

func CreateShortURL(c *fiber.Ctx) error {
	// 同網址不同使用者（包括 guest）產生出來的縮網址不一樣
	// 檢查是 user or guest
	// 使用不同的 get_or_create
	ctx := c.Locals(constants.CTX).(context.Context)
	url := &models.URL{}

	if err := c.BodyParser(url); err != nil {
		return constants.ClientErrorResponse(c, err)
	}

	if err := url.GetOrCreate(ctx); err != nil {
		return constants.ServerErrorResponse(c, err)
	}

	return c.JSON(fiber.Map{
		"short_url": url.ShortURL,
	})
}
