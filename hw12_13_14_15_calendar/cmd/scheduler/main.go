package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/lozhkindm/otus-go-hw/hw12_13_14_15_calendar/internal/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "config", ".env.scheduler", "Path to configuration file")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// loading .env
	if err := godotenv.Load(configFile); err != nil {
		log.Fatal(err)
	}

	// populating config
	config := NewConfig()
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}

	// creating logger
	logg, err := logger.New(config.Logger.Level, config.Logger.Development)
	if err != nil {
		logg.Fatal(err.Error())
	}

	// creating storage
	storage, closeFunc, err := NewStorage(ctx, config)
	if err != nil {
		logg.Fatal("failed to create a storage: " + err.Error())
	}
	defer closeFunc(ctx)

	// creating queue
	queue, err := NewQueue(ctx, config)
	if err != nil {
		logg.Fatal("failed to create a queue client: " + err.Error())
	}
	defer queue.Close(ctx)

	ticker := time.NewTicker(config.App.TickInterval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if err := storage.DeleteOldEvents(ctx); err != nil {
				logg.Error("failed to delete old events: " + err.Error())
			}

			events, err := storage.GetEventsToNotify(ctx)
			if err != nil {
				logg.Error("failed to get events to notify: " + err.Error())
			}

			if len(events) == 0 {
				continue
			}

			for _, event := range events {
				if err := queue.SendEventNotification(ctx, config.RabbitMQ.Exchange.Name, event); err != nil {
					logg.Error("failed to send event notification: " + err.Error())
				} else {
					now := time.Now()
					event.NotifyAt = &now
					if err := storage.UpdateEvent(ctx, event); err != nil {
						logg.Error("failed to set notify at for event: " + err.Error())
					}
				}
			}
		}
	}()

	logg.Info(config.App.Name + " is running...")

	<-ctx.Done()
}
