package main

import (
	"fmt"
	"listener-service/event"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to the rabbitmq
	rabbiConn, err := connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbiConn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbiConn)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.DEBUG", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't contiune until we have a connection
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ, retrying in 1 second...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println("Failed to connect to RabbitMQ after 5 attempts, exiting...")
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
