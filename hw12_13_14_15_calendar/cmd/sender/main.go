package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/lozhkindm/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/lozhkindm/otus-go-hw/hw12_13_14_15_calendar/internal/storage"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "config", ".env.sender", "Path to configuration file")
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

	// creating queue
	queue, err := NewQueue(ctx, config)
	if err != nil {
		logg.Fatal("failed to create a queue client: " + err.Error())
	}
	defer queue.Close(ctx)

	messages, err := queue.Consume(ctx, config.RabbitMQ.Queue.Name, config.RabbitMQ.Consumer.Name)
	if err != nil {
		logg.Fatal("failed to consume messages from a queue: " + err.Error())
	}

	go func() {
		for msg := range messages {
			notification := storage.Notification{}
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				logg.Error("failed to unmarshal message: " + err.Error())
			} else {
				fmt.Println(notification)
			}
		}
	}()

	logg.Info(config.App.Name + " is running...")

	<-ctx.Done()
}
