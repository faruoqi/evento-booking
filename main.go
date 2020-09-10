package main

import (
	"github.com/faruoqi/evento-booking/lib/rabbitmq"
	"github.com/streadway/amqp"
)

const (
	MSGBROKERSERER = "amqp://guest:guest@192.168.17.131:5672"
	QUEUENAME      = "my_queue"
)

func main() {

	connection, err := amqp.Dial(MSGBROKERSERER)
	if err != nil {
		panic(err)
	}

	eventListener, err := rabbitmq.NewAMQPEventListener(connection, QUEUENAME)
	eventListener.Listen("event.created")

}
