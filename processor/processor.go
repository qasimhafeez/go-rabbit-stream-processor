package processor

import (
	"arkis_test/queue"
	"context"

	log "github.com/sirupsen/logrus"
)

type Queue interface {
	Consume(ctx context.Context) (<-chan queue.Delivery, error)
	Publish(ctx context.Context, msg []byte) error
}

type Database interface {
	Get([]byte) (string, error)
}

type processor struct {
	input    Queue
	output   Queue
	database Database
	name     string
}

func New(input, output Queue, db Database, name string) processor {
	return processor{input, output, db, name}
}

func (p processor) Run(ctx context.Context) error {
	deliveries, err := p.input.Consume(ctx)
	if err != nil {
		log.WithError(err).WithField("processor", p.name).Error("Failed to consume messages")
		return err
	}

	log.WithField("processor", p.name).Info("Started consuming messages")

	for {
		select {
		case <-ctx.Done():
			log.WithField("processor", p.name).Info("Processor shutting down")
			return ctx.Err()
		case delivery, ok := <-deliveries:
			if !ok {
				log.WithField("processor", p.name).Info("Message channel closed")
				return nil // Channel is closed, stop the processor
			}
			if err := p.process(ctx, delivery); err != nil {
				log.WithError(err).WithField("processor", p.name).Error("Failed to process delivery")
				return err
			}
		}
	}
}

func (p processor) process(ctx context.Context, delivery queue.Delivery) error {
	inputContent := string(delivery.Body)
	log.WithFields(log.Fields{
		"processor": p.name,
		"delivery":  inputContent,
	}).Info("Processing the delivery")

	data, err := p.database.Get(delivery.Body)
	if err != nil {
		log.WithError(err).WithField("processor", p.name).Error("Error in processing data")
		return err
	}

	log.WithFields(log.Fields{
		"processor": p.name,
		"output":    data,
	}).Info("Publishing processed data to output queue")

	return p.output.Publish(ctx, []byte(data))
}
