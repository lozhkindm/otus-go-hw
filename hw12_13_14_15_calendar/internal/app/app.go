package app

import (
	"context"
	"net/http"

	"github.com/lozhkindm/otus-go-hw/hw12_13_14_15_calendar/internal/storage"

	amqp "github.com/rabbitmq/amqp091-go"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) (int, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID int) error
	DeleteOldEvents(ctx context.Context) error
	ListEvent(ctx context.Context) ([]storage.Event, error)
	GetEventsToNotify(ctx context.Context) ([]storage.Event, error)
	GetEvent(ctx context.Context, eventID int) (*storage.Event, error)
}

type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Queue interface {
	DeclareQueue(ctx context.Context, name string) error
	Consume(ctx context.Context, queueName, consumerName string) (<-chan amqp.Delivery, error)
	Close(ctx context.Context) error
	SendEventNotification(ctx context.Context, queue string, event storage.Event) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) (int, error) {
	// check if date is already booked
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	// check if date is already booked
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, eventID int) error {
	return a.storage.DeleteEvent(ctx, eventID)
}

func (a *App) ListEvent(ctx context.Context) ([]storage.Event, error) {
	return a.storage.ListEvent(ctx)
}

func (a *App) GetEvent(ctx context.Context, eventID int) (*storage.Event, error) {
	return a.storage.GetEvent(ctx, eventID)
}
