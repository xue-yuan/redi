package v1

import (
	"context"
	"fmt"
	"os"
	"redi/config"
	"redi/constants"
	"redi/models"
	"redi/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcuadros/go-defaults"
)

func removeImage(f string) {
	os.Remove(fmt.Sprintf("%s/%s", config.Config.ImageFolder, f))
}

func removeImageByOrder(n, o string) {
	if n != o {
		removeImage(n)
	}
}

func GetShortURLs(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	userID := c.Locals(constants.UserID).(string)

	q := &models.PageQuery{}
	if err := c.QueryParser(q); err != nil {
		return constants.BadRequestResponse(c, err)
	}

	defaults.SetDefaults(q)

	validate := utils.NewValidator()
	if err := validate.Struct(q); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	var urls models.URLPage
	if err := urls.GetAll(ctx, userID, q); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	return constants.OkResponse(c, &fiber.Map{
		"total": urls.Total,
		"rows":  urls.Rows,
	})
}

func GetShortURL(c *fiber.Ctx) error {
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

	url := &models.URL{
		URLID: q.URLID,
	}

	if err := url.Get(ctx, db); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	return constants.OkResponse(c, &fiber.Map{
		"url":         url.URL,
		"url_id":      url.URLID,
		"short_url":   url.ShortURL,
		"created_at":  url.CreatedAt,
		"image":       url.Image,
		"description": url.Description,
		"title":       url.Title,
	})
}

func CreateShortURL(c *fiber.Ctx) error {
	// 同網址不同使用者（包括 guest）產生出來的縮網址不一樣
	// 檢查是 user or guest
	// 使用不同的 get_or_create
	ctx := c.Locals(constants.CTX).(context.Context)
	userID := ""
	url := &models.URL{}

	if err := c.BodyParser(url); err != nil {
		return constants.BadRequestResponse(c, err)
	}

	if s, ok := c.Locals(constants.UserID).(string); ok {
		userID = s
	}

	if err := url.GetOrCreate(ctx, userID); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	return constants.OkResponse(c, &fiber.Map{
		"short_url": url.ShortURL,
	})
}

func DeleteShortURL(c *fiber.Ctx) error {
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
		fmt.Println(0)
		return constants.InternalServerErrorResponse(c, err)
	} else if !ok {
		return constants.ForbiddenResponse(c, constants.ErrOperationNotPermitted)
	}

	tx, err := db.BeginTx(ctx, constants.TxOptions())
	if err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}
	defer tx.Rollback(ctx)

	url := &models.URL{
		URLID: q.URLID,
	}

	if err := url.Get(ctx, tx); err != nil {
		fmt.Println(1)
		return constants.InternalServerErrorResponse(c, err)
	}

	if err := url.Delete(ctx, tx); err != nil {
		fmt.Println(2)
		return constants.InternalServerErrorResponse(c, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	if url.Image != nil {
		removeImage(*url.OpenGraph.Image)
	}

	return constants.EmptyResponse(c)
}

// TODO: 檢查尺寸
func UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	mimeType := file.Header.Get("Content-Type")
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		return constants.BadRequestResponse(c, constants.ErrUnsupportedType)
	}

	fileName := strings.Replace(uuid.New().String(), "-", "", -1)
	fileExt := strings.Split(file.Filename, ".")[1]
	image := fmt.Sprintf("%s.%s", fileName, fileExt)
	dst := fmt.Sprintf("%s/%s", config.Config.ImageFolder, image)

	if err := c.SaveFile(file, dst); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	return constants.OkResponse(c, &fiber.Map{
		"image": image,
	})
}

// 驗證這個 url 是不是該 user 的
// 失敗要把圖片刪掉
func CreateOpenGraph(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	userID := c.Locals(constants.UserID).(string)

	og := &models.OpenGraph{}
	if err := c.BodyParser(og); err != nil {
		removeImage(*og.Image)
		return constants.BadRequestResponse(c, err)
	}

	q := &models.URLIDQuery{}
	if err := c.QueryParser(q); err != nil {
		removeImage(*og.Image)
		return constants.BadRequestResponse(c, err)
	}

	og.URLID = &q.URLID
	validate := utils.NewValidator()
	if err := validate.Struct(og); err != nil {
		removeImage(*og.Image)
		return constants.InternalServerErrorResponse(c, err)
	}

	userURL := &models.UserURL{
		UserID: userID,
		URLID:  q.URLID,
	}

	if ok, err := userURL.HasPermission(ctx); err != nil {
		removeImage(*og.Image)
		return constants.InternalServerErrorResponse(c, err)
	} else if !ok {
		removeImage(*og.Image)
		return constants.ForbiddenResponse(c, constants.ErrOperationNotPermitted)
	}

	if err := og.Create(ctx); err != nil {
		removeImage(*og.Image)
		return constants.InternalServerErrorResponse(c, err)
	}

	return constants.EmptyResponse(c)
}

// 失敗要看新舊圖片有沒有一樣，不一樣要刪掉新的
// 成功要看新舊圖片有沒有一樣，不一樣要刪掉舊的
func UpdateOpenGraph(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	userID := c.Locals(constants.UserID).(string)
	db := ctx.Value(constants.DB).(*pgxpool.Pool)

	og := &models.UpdateOpenGraphRequest{}
	if err := c.BodyParser(og); err != nil {
		removeImageByOrder(og.NewImage, *og.Image)
		return constants.BadRequestResponse(c, err)
	}

	tx, err := db.BeginTx(ctx, constants.TxOptions())
	if err != nil {
		removeImageByOrder(og.NewImage, *og.Image)
		return constants.InternalServerErrorResponse(c, err)
	}
	defer tx.Rollback(ctx)

	q := &models.URLIDQuery{}
	if err := c.QueryParser(q); err != nil {
		removeImageByOrder(og.NewImage, *og.Image)
		return constants.BadRequestResponse(c, err)
	}

	og.URLID = &q.URLID
	validate := utils.NewValidator()
	if err := validate.Struct(og); err != nil {
		removeImageByOrder(og.NewImage, *og.Image)
		return constants.InternalServerErrorResponse(c, err)
	}

	userURL := &models.UserURL{
		UserID: userID,
		URLID:  q.URLID,
	}

	if ok, err := userURL.HasPermission(ctx); err != nil {
		removeImageByOrder(og.NewImage, *og.Image)
		return constants.InternalServerErrorResponse(c, err)
	} else if !ok {
		removeImageByOrder(og.NewImage, *og.Image)
		return constants.ForbiddenResponse(c, constants.ErrOperationNotPermitted)
	}

	if err = og.Update(ctx, tx); err != nil {
		removeImageByOrder(og.NewImage, *og.Image)
		return constants.InternalServerErrorResponse(c, err)
	}

	if err = tx.Commit(ctx); err != nil {
		removeImageByOrder(og.NewImage, *og.Image)
		return constants.InternalServerErrorResponse(c, err)
	}

	removeImageByOrder(*og.Image, og.NewImage)
	return constants.EmptyResponse(c)
}

func DeleteOpenGraph(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	userID := c.Locals(constants.UserID).(string)
	db := ctx.Value(constants.DB).(*pgxpool.Pool)

	q := &models.URLIDQuery{}
	if err := c.QueryParser(q); err != nil {
		return constants.BadRequestResponse(c, err)
	}

	og := &models.OpenGraph{
		URLID: &q.URLID,
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

	if err := og.Delete(ctx, db); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	removeImage(*og.Image)

	return constants.EmptyResponse(c)
}

func CreateCustomizedShortURL(c *fiber.Ctx) error {
	ctx := c.Locals(constants.CTX).(context.Context)
	userID := c.Locals(constants.UserID).(string)
	db := ctx.Value(constants.DB).(*pgxpool.Pool)
	url := &models.URL{}

	if err := c.BodyParser(url); err != nil {
		return constants.BadRequestResponse(c, err)
	}

	if err := url.CreateCustomized(ctx, db, userID); err != nil {
		return constants.InternalServerErrorResponse(c, err)
	}

	return constants.OkResponse(c, &fiber.Map{
		"short_url": url.ShortURL,
	})
}
