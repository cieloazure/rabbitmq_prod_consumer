# Assignment 3

###### Week: 17th-21st September

---
The week from 17th to 21st of September involved discussions on various communication patterns. A brief list of communication patterns described in class are - 
* Remote procedural call
* Message oriented communication using sockets
* ZeroMQ
* Advanced Message Queuing protocol
* Gossip Protocol

Gossip protocol is being implemented as part of the project. For this week's assignment I tried to understand the real world applications of communication pattern using Advanced Message Queuing protocol(AMQP). I used RabbitMQ as a message broker and the golang implementation of AMQP. I demonstrate the publisher-subscriber pattern in this assignment.

### Application

The application is just a crude implementation of dissemination football events to a list of observers. 
The assignment includes two file of importance -
* `football_event_publisher.go`
* `football_event_subscriber.go`

There can be many subscribers and publishers. Each publishing or subscribing their own match. 
### Requirements
* In order to run the program go needs to be installed. Follow the steps given in [here](https://golang.org/doc/install)
* RabbitMQ needs to install. Follow the steps given [here](https://www.rabbitmq.com/download.html)

### Execution
* Start the RabbitMQ server 
```bash
rabbitmq-server
```
* Build the publisher
```bash
go build football_event_publisher.go
```
* Build the subscriber
```bash
go build football_event_subscriber.go
```

* Run the subscriber with following argument
```bash
./football_event_subscriber man_utd arsenal
```

```bash
Usage: ./football_event_subscriber [string] [string]
```

* Run the publisher with the following arguments
```bash
./football_event_publisher man_utd arsenal
```
```bash
Usage: ./football_event_publisher [string] [string]
```

### Possible future extenstion
The following features/specifics are missing but I do intend to add them in the future in order to make a complete project.
- Error Handling
- Test cases
- Read teams from files or an api
- Use probablistic analysis to decide whether an event will occur

