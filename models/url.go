package models

import (
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID        int       `db:"id" json:"id"`
	URLID     uuid.UUID `db:"url_id" json:"url_id" validate:"required,uuid"`
	URL       string    `db:"url" json:"url"`
	ShortURL  string    `db:"short_url" json:"short_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Statistic struct {
	ID         int       `db:"id" json:"id"`
	URLID      uuid.UUID `db:"url_id" json:"url_id" validate:"required,uuid"`
	IPAddress  string    `db:"ip_address" json:"ip_address" validate:"ip"`
	UserAgent  string    `db:"user_agent" json:"user_agent"`
	RefererURL string    `db:"referer_url" json:"referer_url"`
	Latitude   int       `db:"latitude" json:"latitude"`
	Longitude  int       `db:"longitude" json:"longitude"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type OpenGraph struct {
	ID          int       `db:"id" json:"id"`
	URLID       uuid.UUID `db:"url_id" json:"url_id" validate:"required,uuid"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Image       string    `db:"image" json:"image"`
}
