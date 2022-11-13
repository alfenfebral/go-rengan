package pkg_amqp

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

func NewAMQP() (AMQP, error) {
	connection, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		return &AMQPImpl{}, err
	}

	channel, err := connection.Channel()
	if err != nil {
		return &AMQPImpl{}, err
	}

	return &AMQPImpl{
		channel: channel,
	}, err
}

func (a *AMQPImpl) Get() *amqp.Channel {
	return a.channel
}
