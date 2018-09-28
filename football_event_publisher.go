package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
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

func simulateGame(team1 string, team2 string, event_channel string, ch *amqp.Channel) {
	// Set up events
	// TODO: Add more complicated events like offside
	// TODO: Add probablity for each event in order to run simulation
	individual_events := []string{"goal by", "yellow_card for", "red_card for", "foul by"}

	// Set up teams
	teams := make(map[string][]string)

	//TODO: read teams from files or csv
	teams[team1] = []string{"de gea", "shaw", "bailly", "lindelof", "dalot", "pogba", "fred", "matic", "sanchez", "lukaku", "mata"}

	teams[team2] = []string{"cech", "sokratis", "kolasinac", "holding", "xhaka", "elneny", "guendozi", "bellerin", "monreal", "mkhitaryan", "aubameyang"}

	//Simulate game
	timeout := time.After(90 * time.Second)
	tick := time.Tick(1000 * time.Millisecond)
	fmt.Printf("----------------------Starting simulation for %s vs %s----------------\n", team1, team2)

	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			fmt.Println("--------------------------Simulation Complete-----------------------------")
			return
			// Got a tick, we should check on doSomething()
		case <-tick:
			rand.Seed(time.Now().Unix())
			var b strings.Builder
			b.WriteString(individual_events[rand.Intn(len(individual_events))])
			b.WriteString(" ")

			teams_list := []string{team1, team2}
			rand.Seed(time.Now().Unix())
			rand_team := teams_list[rand.Intn(len(teams_list))]

			rand.Seed(time.Now().Unix())
			rand_player := teams[rand_team][rand.Intn(len(teams[rand_team]))]

			b.WriteString(rand_player)
			body := b.String()
			err := ch.Publish(
				event_channel,
				"",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})

			failOnError(err, "Failed to publish message")

			log.Printf("[x]  %s", body)
		}
	}

}

func main() {
	ch := setUpBroker()
	if len(os.Args) != 3 {
		panic(fmt.Sprintf("Usage: ./football_event_publisher [string] [string]"))
	}
	teams := []string{os.Args[1], os.Args[2]}
	event_channel := setUpExchangeForTeams(ch, teams)
	simulateGame(teams[0], teams[1], event_channel, ch)
}
