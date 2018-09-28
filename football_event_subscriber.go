package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func setUpBroker() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return ch
}

func setUpExchangeForTeams(ch *amqp.Channel, teams []string) string {
	teams_string := strings.Join(teams, "_")
	var b strings.Builder
	b.WriteString("events_for_")
	b.WriteString(teams_string)

	err := ch.ExchangeDeclare(
		b.String(),
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare exchange")

	return b.String()
}

func setUpConsumerForEvent(ch *amqp.Channel, event_channel string, team1 string, team2 string) {
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		"",
		event_channel,
		false,
		nil,
	)
	failOnError(err, "Failed to bind to a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("[x]: %s", d.Body)
		}
	}()

	log.Printf("Listening for events of %s vs %s\n", team1, team2)
	<-forever
}

func main() {
	ch := setUpBroker()
	if len(os.Args) != 3 {
		panic(fmt.Sprintf("Usage: ./football_event_subscriber [string] [string]"))
	}
	teams := []string{os.Args[1], os.Args[2]}
	event_channel := setUpExchangeForTeams(ch, teams)
	setUpConsumerForEvent(ch, event_channel, teams[0], teams[1])
}
