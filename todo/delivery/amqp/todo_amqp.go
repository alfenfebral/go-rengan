package todo_amqp_delivery

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
	tracing pkg_tracing.Tracing
	channel pkg_amqp.AMQP
}

// NewTodoAMQPConsumer - make amqp consumer
func NewTodoAMQPConsumer(
	logger pkg_logger.Logger,
	tracing pkg_tracing.Tracing,
	channel pkg_amqp.AMQP,
) TodoAMQPConsumer {
	return &TodoAMQPConsumerImpl{
		logger:  logger,
		tracing: tracing,
		channel: channel,
	}
}

func (consumer *TodoAMQPConsumerImpl) Register() {
	consumer.Create()
}

// Create - create todo consumer
func (consumer *TodoAMQPConsumerImpl) Create() {
	messageName := "send_email"

	channel := consumer.channel.Get()
	q, err := channel.QueueDeclare(messageName, true, false, false, false, nil)
	if err != nil {
		consumer.logger.Error(err)
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
		consumer.logger.Error(err)
	}
	consumer.logger.Println("Consumer listen to queue name", messageName)

	for d := range msgs {
		ctx := pkg_amqp.ExtractAMQPHeaders(context.Background(), d.Headers)

		tr := consumer.tracing.Tracer("amqp")
		opts := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindConsumer),
		}
		_, span := tr.Start(ctx, "AMQP - consume - todo.create", opts...)

		consumer.logger.Printf("Send email to: %s", d.Body)

		err := d.Ack(false)
		if err != nil {
			consumer.logger.Error(err)
		}

		span.End()
	}
}
