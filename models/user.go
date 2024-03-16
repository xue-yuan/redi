package models

import "github.com/google/uuid"

type User struct {
	ID       int       `db:"id" json:"id"`
	UserID   uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
	Username string    `db:"username" json:"username"`
	Password string    `db:"password" json:"password"`
}

type UserURL struct {
	ID     int       `db:"id" json:"id"`
	URLID  uuid.UUID `db:"url_id" json:"url_id" validate:"required,uuid"`
	UserID uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
}
