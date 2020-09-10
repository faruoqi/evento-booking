package rabbitmq

import (
	"github.com/faruoqi/evento/contracts"
	"github.com/faruoqi/evento/messaging"
	"log"
)

type EventProcessor struct {
	EventListener messaging.EventListener
}

func (e *EventProcessor) ProcessEvents() error {

	log.Print("Listening to Events")

	received, errors, err := e.EventListener.Listen("event.created")
	if err != nil {
		return err
	}
	for {
		select {
		case evt := <-received:
			e.HandleEvent(evt)
		case err = <-errors:
			log.Printf("received error when processing message: %s", err)
		}
	}
}

func (e *EventProcessor) HandleEvent(event messaging.Event) {

	switch e := event.(type) {
	case *contracts.EventCreatedEvent:
		log.Printf("event %s at %s created %s", e.Name, e.LocationName, e)
	default:
		log.Printf("unknown event %t", e)
	}

}
