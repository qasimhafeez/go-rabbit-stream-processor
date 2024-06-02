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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Queues for Stream A
	inputQueueA, err := queue.New(os.Getenv("RABBITMQ_URL"), "input-A")
	if err != nil {
		log.WithError(err).Panic("Cannot create input queue A")
	}
	outputQueueA, err := queue.New(os.Getenv("RABBITMQ_URL"), "output-A")
	if err != nil {
		log.WithError(err).Panic("Cannot create output queue A")
	}

	// Queues for Stream B
	inputQueueB, err := queue.New(os.Getenv("RABBITMQ_URL"), "input-B")
	if err != nil {
		log.WithError(err).Panic("Cannot create input queue B")
	}
	outputQueueB, err := queue.New(os.Getenv("RABBITMQ_URL"), "output-B")
	if err != nil {
		log.WithError(err).Panic("Cannot create output queue B")
	}

	log.Info("Application is ready to run")

	go processor.New(inputQueueA, outputQueueA, database.D{}).Run(ctx)
	go processor.New(inputQueueB, outputQueueB, database.D{}).Run(ctx)

	select {}
}
