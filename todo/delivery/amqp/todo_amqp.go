package todo_amqp

import (
	"context"

	pkg_amqp "go-rengan/pkg/amqp"
	pkg_logger "go-rengan/pkg/logger"
	pkg_tracing "go-rengan/pkg/tracing"

	"go.opentelemetry.io/otel/trace"
)

type TodoAMQPConsumer interface {
	Create()
	Register()
}

// TodoAMQPConsumerImpl represent the amqp
type TodoAMQPConsumerImpl struct {
	logger  pkg_logger.Logger
	tp      pkg_tracing.Tracing
	channel pkg_amqp.AMQP
}

// NewTodoAMQPConsumer - make amqp consumer
func NewTodoAMQPConsumer(
	logger pkg_logger.Logger,
	tp pkg_tracing.Tracing,
	channel pkg_amqp.AMQP,
) TodoAMQPConsumer {
	return &TodoAMQPConsumerImpl{
		logger:  logger,
		tp:      tp,
		channel: channel,
	}
}

func (c *TodoAMQPConsumerImpl) Register() {
	c.Create()
}

// Create - create todo consumer
func (c *TodoAMQPConsumerImpl) Create() {
	messageName := "todo.create"

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
		ctx := pkg_amqp.ExtractAMQPHeaders(context.Background(), d.Headers)

		tr := c.tp.Tracer("amqp")
		opts := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindConsumer),
		}
		_, span := tr.Start(ctx, "AMQP - consume - todo.create", opts...)

		c.logger.Printf("Received a message: %s", d.Body)

		err := d.Ack(false)
		if err != nil {
			c.logger.Error(err)
		}

		span.End()
	}
}
