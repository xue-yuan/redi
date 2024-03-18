package models

import (
	"context"
	"errors"
	"redi/constants"
	"redi/utils"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const MAX_ATTEMPT = 10

type URL struct {
	ID        int       `db:"id" json:"id"`
	URLID     string    `db:"url_id" json:"url_id"`
	URL       string    `db:"url" json:"url"`
	ShortURL  string    `db:"short_url" json:"short_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Statistic struct {
	ID         int       `db:"id" json:"id"`
	URLID      string    `db:"url_id" json:"url_id"`
	IPAddress  string    `db:"ip_address" json:"ip_address" validate:"ip"`
	UserAgent  string    `db:"user_agent" json:"user_agent"`
	RefererURL string    `db:"referer_url" json:"referer_url"`
	Latitude   int       `db:"latitude" json:"latitude"`
	Longitude  int       `db:"longitude" json:"longitude"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type OpenGraph struct {
	ID          int    `db:"id" json:"id"`
	URLID       string `db:"url_id" json:"url_id" validate:"required,uuid"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Image       string `db:"image" json:"image"`
}

// for guest
// HACK: what if error occured during rollback and commit?
func (u *URL) GetOrCreate(ctx context.Context) error {
	db := ctx.Value(constants.DB).(*pgxpool.Pool)
	tx, err := db.BeginTx(ctx, constants.TxOptions())
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	getQuery := `
		SELECT u.id, u.url, u.url_id, u.short_url, u.created_at
		FROM urls AS u
		LEFT JOIN user_urls AS uu ON u.url_id = uu.url_id
		WHERE uu.url_id IS NULL AND u.url = $1;
	`

	rows, err := tx.Query(ctx, getQuery, u.URL)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&u.ID, &u.URL, &u.URLID, &u.ShortURL, &u.CreatedAt); err != nil {
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
		return nil
	}

	for i := 0; i < MAX_ATTEMPT; i++ {
		u.ShortURL = utils.GenerateShortURL(seed, i)
		tag, err = tx.Exec(ctx, insertQuery, u.URLID, u.URL, u.ShortURL)
		if err != nil {
			return err
		} else if tag.RowsAffected() > 0 {
			break
		}

		seed = u.ShortURL
	}

	if tag.RowsAffected() < 1 {
		return errors.New("create failed")
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// for user
func GetOrCreateByUser() {

}
