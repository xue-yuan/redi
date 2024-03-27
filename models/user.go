package models

import (
	"redi/constants"
	"redi/utils"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

type User struct {
	ID        int       `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UserURL struct {
	ID     int    `db:"id" json:"id"`
	UserID string `db:"user_id" json:"user_id"`
	URLID  string `db:"url_id" json:"url_id"`
}

func (u *User) Login(ctx context.Context, p string) (bool, error) {
	db := ctx.Value(constants.DB).(*pgxpool.Pool)
	query := `
		SELECT *
		FROM users
		WHERE username = $1;
	`

	rows, err := db.Query(ctx, query, u.Username)
	if err != nil {
		return false, err
	}

	*u, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return utils.IsValidPassword(u.Password, p), nil
}

func (u *User) Create(ctx context.Context) error {
	var err error
	db := ctx.Value(constants.DB).(*pgxpool.Pool)
	query := `
		INSERT INTO users (user_id, username, password)
		VALUES ($1, $2, $3)
		RETURNING *
	`

	u.UserID, err = utils.GenerateUUID("user")
	if err != nil {
		return err
	}

	h, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	rows, err := db.Query(ctx, query, u.UserID, u.Username, h)
	if err != nil {
		return err
	}

	*u, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return constants.ErrDuplicateUsername
		}

		return err
	}

	return nil
}

func (u *UserURL) HasPermission(ctx context.Context) (bool, error) {
	db := ctx.Value(constants.DB).(*pgxpool.Pool)
	query := `
		SELECT COUNT(*)
		FROM user_urls
		WHERE user_id = $1 AND url_id = $2
	`
	count := 0

	if err := db.QueryRow(ctx, query, u.UserID, u.URLID).Scan(&count); err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
