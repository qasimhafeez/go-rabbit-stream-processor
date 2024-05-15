package main

import (
	"arkis_test/database"
	"arkis_test/processor"
	"arkis_test/queue"
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	inputQueue, err := queue.New(os.Getenv("RABBITMQ_URL"), "input-A")
	if err != nil {
		log.WithError(err).Panic("Cannot create input queue")
	}

	outputQueue, err := queue.New(os.Getenv("RABBITMQ_URL"), "output-A")
	if err != nil {
		log.WithError(err).Panic("Cannot create output queue")
	}

	log.Info("Application is ready to run")

	processor.New(inputQueue, outputQueue, database.D{}).Run(ctx)
}
