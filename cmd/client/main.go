package main

import (
	"fmt"
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
	gameState := gamelogic.NewGameState(username)

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	err = pubsub.SubscribeToJSON(connection, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient, handlerPause(gameState))
	if err != nil {
		fmt.Println("❌ Subscription error:", err)
		return
	}
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
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(state routing.PlayingState) {
		fmt.Println("Sate: ", state)
		defer fmt.Print("> ")
		gs.HandlePause(state)
	}
}
