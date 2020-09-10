package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/faruoqi/evento/contracts"
	"github.com/faruoqi/evento/messaging"
	"github.com/streadway/amqp"
	"log"
)

type amqpEventListener struct {
	connection *amqp.Connection
	queue      string
}

func NewAMQPEventListener(conn *amqp.Connection, queue string) (*amqpEventListener, error) {

	listener := &amqpEventListener{
		connection: conn,
		queue:      queue,
	}

	err := listener.setup()
	if err != nil {
		return nil, err
	}
	return listener, nil

}

func (a *amqpEventListener) setup() error {

	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	_, err = channel.QueueDeclare(a.queue, true, false, false, false, nil)
	return err
}

func (a *amqpEventListener) Listen(eventNames ...string) {

	channel, err := a.connection.Channel()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer channel.Close()
	for _, eventName := range eventNames {
		if err := channel.QueueBind(a.queue, eventName, "events", false, nil); err != nil {
			log.Fatal(err.Error())
		}
	}

	msgs, err := channel.Consume(a.queue, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	for msg := range msgs {
		rawEventName, ok := msg.Headers["x-event-name"]
		if !ok {
			log.Print(fmt.Errorf("msg did not contain x-event-name header"))
			msg.Nack(false, false)
			continue
		}
		eventName, ok := rawEventName.(string)
		if !ok {
			log.Print(fmt.Errorf("x-event-name header is not string , but %t", eventName))
			continue
		}
		var event messaging.Event
		switch eventName {
		case "event.created":
			event = new(contracts.EventCreatedEvent)
		default:
			log.Println(fmt.Errorf("event type %s unknown ", eventName))
		}
		err := json.Unmarshal(msg.Body, event)
		if err != nil {
			log.Print(err.Error())
			continue
		}

		switch e := event.(type) {
		case *contracts.EventCreatedEvent:
			log.Printf("event %s at %s created ", e.Name, e.LocationName)
		default:
			log.Printf("unknown event %t", e)
		}
		msg.Ack(false)

	}

}
