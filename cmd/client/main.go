package main

import (
	"fmt"
	"os"
	"os/signal"

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

	fmt.Println("✅ Queue successfully declared and bound.")

	osChan := make(chan os.Signal)
	signal.Notify(osChan, os.Interrupt)
	<-osChan
}
