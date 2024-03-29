package models

import (
	"context"
	"redi/database"
	"time"

	"github.com/jackc/pgx/v5"
)

type Statistic struct {
	ID          int       `db:"int" json:"id"`
	URLID       string    `db:"url_id" json:"url_id"`
	IPAddress   string    `db:"ip_address" json:"ip_address" validate:"ip"`
	UserAgent   string    `db:"user_agent" json:"user_agent"`
	RefererURL  string    `db:"referer_url" json:"referer_url"`
	CountryCode string    `db:"country_code" json:"country_code"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type Statistics []Statistic

func (s *Statistic) Count(ctx context.Context, tx database.Tx) (count int, err error) {
	query := `
		SELECT COUNT(*)
		FROM statistics
		WHERE url_id = $1
	`

	if err = tx.QueryRow(ctx, query, s.URLID).Scan(&count); err != nil {
		return
	}

	return
}

func (s *Statistics) GetAll(ctx context.Context, tx database.Tx, urlID string) (err error) {
	query := `
		SELECT *
		FROM statistics
		WHERE url_id = $1
	`

	rows, err := tx.Query(ctx, query, urlID)
	if err != nil {
		return
	}

	*s, err = pgx.CollectRows(rows, pgx.RowToStructByName[Statistic])
	if err != nil {
		return
	}

	return
}

func (s *Statistic) Create(ctx context.Context, tx database.Tx) (err error) {
	query := `
		INSERT INTO statistics (url_id, ip_address, user_agent, referer_url, country_code)
		VALUES ($1, $2, $3, $4, $5)
	`

	if _, err = tx.Exec(
		ctx,
		query,
		s.URLID,
		s.IPAddress,
		s.UserAgent,
		s.RefererURL,
		s.CountryCode,
	); err != nil {
		return err
	}

	return nil
}
