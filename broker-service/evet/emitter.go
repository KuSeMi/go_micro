package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Publishing event:", event, "with severity:", severity)

	body := []byte(event)
	err = channel.Publish(
		"logs_topic", // exchange
		severity,     // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(con *amqp.Connection) (*Emitter, error) {
	emitter := &Emitter{
		conn: con,
	}

	err := emitter.setup()
	if err != nil {
		return nil, err
	}

	return emitter, nil
}
