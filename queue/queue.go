package queue

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Delivery struct {
	Body    []byte
	handler amqp.Delivery
}

type queue struct {
	amqpConnection *amqp.Connection
	channel        *amqp.Channel
	name           string
}

func New(connectionURL, queueName string) (*queue, error) {
	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	return &queue{conn, ch, queueName}, nil
}

func (queue *queue) Consume(ctx context.Context) (<-chan Delivery, error) {
	deliveries, err := queue.channel.ConsumeWithContext(
		ctx,
		queue.name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	out := make(chan Delivery)

	go func() {
		for {
			select {
			case delivery := <-deliveries:
				out <- Delivery{delivery.Body, delivery}
			case <-ctx.Done():
				close(out)
				return
			}
		}
	}()

	return out, nil
}

func (queue *queue) Publish(ctx context.Context, msg []byte) error {
	data := amqp.Publishing{
		DeliveryMode:    amqp.Transient,
		Timestamp:       time.Now(),
		Body:            msg,
		ContentEncoding: "application/json",
	}

	return queue.channel.PublishWithContext(ctx, "", queue.name, true, false, data)
}
