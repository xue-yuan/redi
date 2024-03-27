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
		return constants.NotFoundResponse(c)
	} else if err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	stat := &models.Statistic{
		URLID: u.URLID,
	}

	if err := stat.Create(ctx, db); err != nil {
		// log
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
