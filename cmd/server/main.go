package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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
	fmt.Println("Connection was successful")

loop:
	for {
		gamelogic.PrintServerHelp()
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}
		fmt.Println(inputs)
		switch strings.Join(inputs, "") {
		case "pause":
			pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})

		case "resume":
			pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
		case "quit":
			break loop
		default:
			fmt.Println("Did not understand the command")
			continue
		}
	}

	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	<-osChan
}
