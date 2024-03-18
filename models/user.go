package models

import (
	"errors"
	"redi/constants"
	"redi/utils"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	URLID  string `db:"url_id" json:"url_id"`
	UserID string `db:"user_id" json:"user_id"`
}

func hashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func IsValidPassword(p, h string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
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

	*u, err = pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName[User])
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return IsValidPassword(p, u.Password), nil
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

	h, err := hashPassword(u.Password)
	if err != nil {
		return err
	}

	rows, err := db.Query(ctx, query, u.UserID, u.Username, h)
	if err != nil {
		return err
	}

	*u, err = pgx.CollectExactlyOneRow[User](rows, pgx.RowToStructByName[User])
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return errors.New("duplicate username")
		}

		return err
	}

	return nil
}
