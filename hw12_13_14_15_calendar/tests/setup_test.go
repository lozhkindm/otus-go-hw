package tests

import (
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	db   *sqlx.DB
	host string
)

type eventCreateRequest struct {
	UserID      int     `json:"user_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	StartAt     int64   `json:"start_at"`
	EndAt       int64   `json:"end_at"`
	NotifyAt    *int64  `json:"notify_at"`
}

type eventUpdateRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	StartAt     int64   `json:"start_at"`
	EndAt       int64   `json:"end_at"`
	NotifyAt    *int64  `json:"notify_at"`
}

type eventResponse struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       time.Time  `json:"end_at"`
	NotifyAt    *time.Time `json:"notify_at"`
}

type eventListResponse struct {
	Events []eventResponse `json:"events"`
}

func TestMain(m *testing.M) {
	var err error
	time.Local = nil
	db, err = sqlx.Open("pgx", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal("failed to open pgx")
	}
	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping db")
	}
	host = os.Getenv("HTTP_HOST")
	os.Exit(m.Run())
}
