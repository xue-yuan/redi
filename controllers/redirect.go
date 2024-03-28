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
		ShortURL:  shortURL,
		OpenGraph: models.OpenGraph{},
	}

	if err := u.HMGetOpenGraphByShortURL(ctx); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	if u.URL == "" {
		if err := u.GetOpenGraphByShortURL(ctx, db); err == pgx.ErrNoRows {
			return constants.NotFoundResponse(c)
		} else if err != nil {
			return constants.InternalServerErrorResponse(c, err)
		}

		if err := u.HMSetOpenGraphByShortURL(ctx); err != nil {
			return constants.InternalServerErrorResponse(c, err)
		}
	}

	stat := &models.Statistic{
		URLID:      u.URLID,
		IPAddress:  utils.GetIP(c),
		UserAgent:  c.Get("User-Agent"),
		RefererURL: c.Get("Referer"),
	}

	if err := stat.Create(ctx, db); err != nil {
		// log
		fmt.Println(err)
	}

	if utils.IsStructEmpty(u.OpenGraph) {
		return constants.RedirectResponse(c, u.URL)
	}

	return c.Render("index", fiber.Map{
		"URL":         u.URL,
		"ShortURL":    u.ShortURL,
		"Title":       u.Title,
		"Description": u.Description,
		"Image":       fmt.Sprintf("/image/%s", u.Image),
	})
}
