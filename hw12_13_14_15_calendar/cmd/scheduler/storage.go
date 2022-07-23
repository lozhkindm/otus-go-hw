package main

import (
	"context"
	"errors"

	"github.com/lozhkindm/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/lozhkindm/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/lozhkindm/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	storageTypeMemory = "memory"
	storageTypeSQL    = "sql"
)

var undefinedStorageType = errors.New("undefined storage type")

func NewStorage(ctx context.Context, config Config) (app.Storage, func(context.Context) error, error) {
	var (
		storage   app.Storage
		closeFunc func(ctx2 context.Context) error
	)

	switch config.App.StorageType {
	case storageTypeMemory:
		storage = memorystorage.New()
		closeFunc = func(ctx context.Context) error {
			return nil
		}
	case storageTypeSQL:
		db, err := sqlx.Open("pgx", config.PostgreSQL.BuildDSN())
		if err != nil {
			return nil, nil, err
		}
		ss := sqlstorage.New(db)
		if err := ss.Connect(ctx); err != nil {
			return nil, nil, err
		}
		storage = ss
		closeFunc = ss.Close
	default:
		return nil, nil, undefinedStorageType
	}

	return storage, closeFunc, nil
}
