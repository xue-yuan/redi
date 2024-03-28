package models

import (
	"context"
	"fmt"
	"redi/config"
	"redi/constants"
	"redi/database"
	"redi/redis"
	"redi/utils"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const MAX_ATTEMPT = 10

type OpenGraph struct {
	URLID       string `db:"url_id" json:"url_id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Image       string `db:"image" json:"image" validate:"exist"`
}

type URL struct {
	OpenGraph
	URLID     string    `db:"url_id" json:"url_id"`
	URL       string    `db:"url" json:"url"`
	ShortURL  string    `db:"short_url" json:"short_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type URLs []URL

type URLPage struct {
	Total int  `db:"total" json:"total"`
	Rows  URLs `db:"rows" json:"rows"`
}

func (u *URLPage) GetAll(ctx context.Context, userID string, q *PageQuery) error {
	db := ctx.Value(constants.DB).(*pgxpool.Pool)
	query := fmt.Sprintf(`
		WITH total_rows AS (
			SELECT COUNT(*) AS total
			FROM urls AS u
			LEFT JOIN user_urls AS uu ON u.url_id = uu.url_id
			WHERE user_id = $1
		),
		rows_data AS (
			SELECT u.*
			FROM urls AS u
			LEFT JOIN user_urls AS uu ON u.url_id = uu.url_id
			WHERE user_id = $1
			ORDER BY url %s
			LIMIT $2 OFFSET $3
		)
		SELECT
			json_build_object(
				'total', (SELECT total FROM total_rows),
				'rows', json_agg(rd)
			) AS result
		FROM rows_data rd;
	`, q.Order)

	if err := db.QueryRow(ctx, query, userID, q.Limit, q.Offset).Scan(u); err != nil {
		return err
	}

	return nil
}

// func NewOpenGraph() *OpenGraph {
// 	return &OpenGraph{
// 		Title:       "",
// 		Description: "",
// 		Image:       "",
// 	}
// }

func (u *URL) Get(ctx context.Context, tx database.Tx) error {
	query := `
		SELECT u.url, u.url_id, u.short_url, u.created_at, o.title, o.description, o.image
		FROM urls AS u
		LEFT JOIN open_graphs AS o ON u.url_id = o.url_id
		WHERE u.url_id = $1
	`

	rows, err := tx.Query(ctx, query, u.URLID)
	if err != nil {
		return err
	}

	*u, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[URL])
	if err != nil {
		return err
	}

	return nil
}

func (u *URL) GetOrCreate(ctx context.Context, userID string) error {
	db := ctx.Value(constants.DB).(*pgxpool.Pool)
	tx, err := db.BeginTx(ctx, constants.TxOptions())
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	args := []interface{}{u.URL}
	getQuery := ""
	if userID == "" {
		getQuery = `
			SELECT u.url, u.url_id, u.short_url, u.created_at
			FROM urls AS u
			LEFT JOIN user_urls AS uu ON u.url_id = uu.url_id
			WHERE uu.url_id IS NULL AND u.url = $1;
		`
	} else {
		args = append(args, userID)
		getQuery = `
			SELECT u.url, u.url_id, u.short_url, u.created_at
			FROM urls AS u
			LEFT JOIN user_urls AS uu ON u.url_id = uu.url_id
			WHERE uu.url_id IS NOT NULL AND u.url = $1 AND uu.user_id = $2
		`
	}

	rows, err := tx.Query(ctx, getQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&u.URL, &u.URLID, &u.ShortURL, &u.CreatedAt); err != nil {
			return err
		}

		return nil
	}

	var tag pgconn.CommandTag
	var seed string = u.URL
	insertQuery := `
		INSERT INTO urls (url_id, url, short_url)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING;
	`

	u.URLID, err = utils.GenerateUUID("url")
	if err != nil {
		return err
	}

	for i := 0; i < MAX_ATTEMPT; i++ {
		u.ShortURL = utils.GenerateShortURL(seed, userID, i)
		tag, err = tx.Exec(ctx, insertQuery, u.URLID, u.URL, u.ShortURL)
		if err != nil {
			return err
		} else if tag.RowsAffected() > 0 {
			break
		}

		seed = u.ShortURL
	}

	if tag.RowsAffected() < 1 {
		return constants.ErrCreateShortURLFailed
	}

	if userID != "" {
		insertQuery = `
			INSERT INTO user_urls (url_id, user_id)
			VALUES ($1, $2)
		`
		if _, err := tx.Exec(ctx, insertQuery, u.URLID, userID); err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (u *URL) CreateCustomized(ctx context.Context, tx database.Tx, userID string) error {
	var err error
	query := `
		INSERT INTO urls (url_id, url, short_url)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING;
	`

	u.URLID, err = utils.GenerateUUID("url")
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, query, u.URLID, u.URL, u.ShortURL); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return constants.ErrDuplicateShortURL
		}

		return err
	}

	userQuery := `
		INSERT INTO user_urls (url_id, user_id)
		VALUES ($1, $2)
	`

	if _, err := tx.Exec(ctx, userQuery, u.URLID, userID); err != nil {
		return err
	}

	return nil
}

func (u *URL) Delete(ctx context.Context, tx database.Tx) error {
	query := `
		DELETE FROM urls
		WHERE url_id = $1
	`

	if _, err := tx.Exec(ctx, query, u.URLID); err != nil {
		return err
	}

	return nil
}

func (u *URL) GetOpenGraphByShortURL(ctx context.Context, tx database.Tx) error {
	query := `
		SELECT u.url, u.url_id, u.short_url, u.created_at, o.title, o.description, o.image
		FROM urls AS u
		LEFT JOIN open_graphs AS o ON o.url_id = u.url_id
		WHERE u.short_url = $1
	`

	rows, err := tx.Query(ctx, query, u.ShortURL)
	if err != nil {
		return err
	}

	if *u, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[URL]); err != nil {
		return err
	}

	return nil
}

func (u *URL) HMGetOpenGraphByShortURL(ctx context.Context) error {
	fields := []string{"url", "url_id", "title", "description", "image"}

	result, err := redis.Client.HMGet(ctx, u.ShortURL, fields...).Result()
	if err != nil {
		return err
	}

	for i, field := range fields {
		if result[i] != nil {
			switch field {
			case "url":
				u.URL = result[i].(string)
			case "url_id":
				u.URLID = result[i].(string)
			case "title":
				u.Title = result[i].(string)
			case "description":
				u.Description = result[i].(string)
			case "image":
				u.Image = result[i].(string)
			}
		}
	}

	return nil
}

func (u *URL) HMSetOpenGraphByShortURL(ctx context.Context) error {
	if err := redis.Client.HMSet(ctx, u.ShortURL, map[string]interface{}{
		"url":         u.URL,
		"url_id":      u.URLID,
		"title":       u.Title,
		"description": u.Description,
		"image":       u.Image,
	}).Err(); err != nil {
		return err
	}

	if err := redis.Client.Expire(ctx, u.ShortURL, config.Config.ShortURLTTL).Err(); err != nil {
		return err
	}

	return nil
}

func (o *OpenGraph) Create(ctx context.Context) error {
	db := ctx.Value(constants.DB).(*pgxpool.Pool)
	query := `
		INSERT INTO open_graphs (url_id, title, description, image)
		VALUES ($1, $2, $3, $4)
	`

	if _, err := db.Exec(ctx, query, o.URLID, o.Title, o.Description, o.Image); err != nil {
		return err
	}

	return nil
}

func (o *OpenGraph) Update(ctx context.Context, tx database.Tx) error {
	query := `
		UPDATE open_graphs
		SET title = $1, description = $2, image = $3
		WHERE url_id = $4
	`

	if _, err := tx.Exec(ctx, query, o.Title, o.Description, o.Image, o.URLID); err != nil {
		return err
	}

	return nil
}

func (o *OpenGraph) Delete(ctx context.Context, tx database.Tx) error {
	getQuery := `
		SELECT image
		FROM open_graphs
		WHERE url_id = $1
	`

	if err := tx.QueryRow(ctx, getQuery, o.URLID).Scan(&o.Image); err != nil {
		return err
	}

	deleteQuery := `
		DELETE FROM open_graphs
		WHERE url_id = $1
	`

	if _, err := tx.Exec(ctx, deleteQuery, o.URLID); err != nil {
		return err
	}

	return nil
}

func (o *OpenGraph) GetImage(ctx context.Context, tx pgx.Tx) (string, error) {
	query := `
		SELECT image
		FROM open_graphs
		WHERE url_id = $1
	`

	oldImage := ""
	if err := tx.QueryRow(ctx, query, o.URLID).Scan(&oldImage); err != nil {
		return oldImage, err
	}

	return oldImage, nil
}
