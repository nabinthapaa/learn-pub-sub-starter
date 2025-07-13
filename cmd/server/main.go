package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const CONNECTION_STRING = "amqp://guest:guest@localhost:5672"

func main() {
	fmt.Println("Starting Peril server...")
	connection, error := amqp.Dial(CONNECTION_STRING)
	if error != nil {
		fmt.Println("Error creating connection")
		return
	}
	defer connection.Close()
	channel, _ := connection.Channel()
	pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	fmt.Println("Connection was successful")

	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	<-osChan
}
