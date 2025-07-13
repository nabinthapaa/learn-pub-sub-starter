package main

import (
	"fmt"
	"os"
	"os/signal"
	"slices"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const CONNECTION_STRING = "amqp://guest:guest@localhost:5672"

func main() {
	fmt.Println("Starting Peril client...")

	connection, err := amqp.Dial(CONNECTION_STRING)
	if err != nil {
		fmt.Println("Something went wrong while creating connection")
		return
	}
	defer connection.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("Something went wrong with user input")
		return
	}
	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	_, _, err = pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient)
	if err != nil {
		fmt.Println("❌ Queue declaration/binding error:", err)
		return
	}

	gameState := gamelogic.NewGameState(username)
	gamelogic.PrintClientHelp()
loop:
	for {
		input := gamelogic.GetInput()

		switch true {
		case slices.Contains(input, "move"):
			_, err = gameState.CommandMove(input)
			if err != nil {
				fmt.Println("Command not allowed")
				continue
			}

		case slices.Contains(input, "spawn"):
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println("Command not allowed")
				continue
			}

		case slices.Contains(input, "status"):
			gameState.CommandStatus()

		case slices.Contains(input, "help"):
			gamelogic.PrintClientHelp()

		case slices.Contains(input, "spam"):
			fmt.Println("Spamming not allowed yet!")

		case slices.Contains(input, "quit"):
			break loop

		default:
			fmt.Println("Did not understand the command")
			continue

		}
	}

	fmt.Println("✅ Queue successfully declared and bound.")

	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	<-osChan
}
