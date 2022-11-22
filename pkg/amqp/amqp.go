package amqp

import (
	"os"

	"github.com/streadway/amqp"
)

type AMQP interface {
	Get() *amqp.Channel
}

type AMQPImpl struct {
	channel *amqp.Channel
}

func New() (AMQP, error) {
	connection, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	return &AMQPImpl{
		channel: channel,
	}, err
}

func (a *AMQPImpl) Get() *amqp.Channel {
	return a.channel
}
