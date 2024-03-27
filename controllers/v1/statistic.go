package v1

import (
	"context"
	"redi/constants"
	"redi/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StatCount(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	userID := c.Locals(constants.UserID).(string)
	db := ctx.Value(constants.DB).(*pgxpool.Pool)

	q := &models.URLIDQuery{}
	if err := c.QueryParser(q); err != nil {
		return constants.BadRequestResponse(c, err)
	}

	userURL := &models.UserURL{
		UserID: userID,
		URLID:  q.URLID,
	}

	if ok, err := userURL.HasPermission(ctx); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	} else if !ok {
		return constants.ForbiddenResponse(c, constants.ErrOperationNotPermitted)
	}

	stat := &models.Statistic{
		URLID: q.URLID,
	}

	count, err := stat.Count(ctx, db)
	if err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	return constants.OkResponse(c, &fiber.Map{
		"count": count,
	})
}

func GetStats(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	userID := c.Locals(constants.UserID).(string)
	db := ctx.Value(constants.DB).(*pgxpool.Pool)

	q := &models.URLIDQuery{}
	if err := c.QueryParser(q); err != nil {
		return constants.BadRequestResponse(c, err)
	}

	userURL := &models.UserURL{
		UserID: userID,
		URLID:  q.URLID,
	}

	if ok, err := userURL.HasPermission(ctx); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	} else if !ok {
		return constants.ForbiddenResponse(c, constants.ErrOperationNotPermitted)
	}

	stats := &models.Statistics{}
	if err := stats.GetAll(ctx, db, q.URLID); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	return constants.OkResponse(c, &fiber.Map{
		"data": stats,
	})
}
