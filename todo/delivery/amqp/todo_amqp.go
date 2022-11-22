package amqpdelivery

import (
	"context"

	pkgamqp "go-rengan/pkg/amqp"
	logger "go-rengan/pkg/logger"
	tracing "go-rengan/pkg/tracing"

	"go.opentelemetry.io/otel/trace"
)

type AMQPConsumer interface {
	Create()
	Register()
}

// AMQPConsumerImpl represent the amqp
type AMQPConsumerImpl struct {
	logger  logger.Logger
	tracing tracing.Tracing
	channel pkgamqp.AMQP
}

// New - make amqp consumer
func New(
	logger logger.Logger,
	tracing tracing.Tracing,
	channel pkgamqp.AMQP,
) AMQPConsumer {
	return &AMQPConsumerImpl{
		logger:  logger,
		tracing: tracing,
		channel: channel,
	}
}

func (c *AMQPConsumerImpl) Register() {
	c.Create()
}

// Create - create todo consumer
func (c *AMQPConsumerImpl) Create() {
	messageName := "send_email"

	channel := c.channel.Get()
	q, err := channel.QueueDeclare(messageName, true, false, false, false, nil)
	if err != nil {
		c.logger.Error(err)
	}

	msgs, err := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Error(err)
	}
	c.logger.Println("Consumer listen to queue name", messageName)

	for d := range msgs {
		ctx := pkgamqp.ExtractAMQPHeaders(context.Background(), d.Headers)

		tr := c.tracing.Tracer("amqp")
		opts := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindConsumer),
		}
		_, span := tr.Start(ctx, "AMQP - consume - todo.create", opts...)

		c.logger.Printf("Send email to: %s", d.Body)

		err := d.Ack(false)
		if err != nil {
			c.logger.Error(err)
		}

		span.End()
	}
}
