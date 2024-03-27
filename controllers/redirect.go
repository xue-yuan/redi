package controllers

import (
	"context"
	"fmt"
	"redi/constants"
	"redi/models"
	"redi/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RedirectURL(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	db := ctx.Value(constants.DB).(*pgxpool.Pool)
	shortURL := c.Params("short_url")

	u := &models.URL{
		ShortURL: shortURL,
	}

	if err := u.GetOpenGraphByShortURL(ctx, db); err == pgx.ErrNoRows {
		fmt.Println(1, err)
		return constants.NotFoundResponse(c)
	} else if err != nil {
		fmt.Println(2, err)
		return constants.InternalServerErrorResponse(c, err)
	}

	if utils.IsStructEmpty(u.OpenGraph) {
		return constants.RedirectResponse(c, u.URL)
	}

	return c.Render("index", fiber.Map{
		"URL":         u.URL,
		"ShortURL":    u.ShortURL,
		"Title":       u.Title,
		"Description": u.Description,
		"Image":       fmt.Sprintf("image/%s", *u.Image),
	})
}
