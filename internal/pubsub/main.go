package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	Transient SimpleQueueType = "transient"
	Durable   SimpleQueueType = "Durable"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	data, error := json.Marshal(val)
	if error != nil {
		return error
	}

	error = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})

	if error != nil {
		return error
	}

	return nil
}

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	isDurable := queueType == Durable
	isTransient := queueType == Transient

	fmt.Println(isDurable, isTransient, queueName)

	queue, err := channel.QueueDeclare(queueName, isDurable, isTransient, isTransient, false, nil)
	if err != nil {
		return nil, queue, err
	}
	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, queue, err
	}

	return channel, queue, nil
}
