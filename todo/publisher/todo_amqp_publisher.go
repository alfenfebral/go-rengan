package amqppublisher

import (
	"context"
	"fmt"

	pkgamqp "go-rengan/pkg/amqp"
	logger "go-rengan/pkg/logger"
	tracing "go-rengan/pkg/tracing"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

type AMQPPublisher interface {
	Create(value string)
}

type AMQPPublisherImpl struct {
	logger  logger.Logger
	tracing tracing.Tracing
	channel pkgamqp.AMQP
}

func New(
	logger logger.Logger,
	tracing tracing.Tracing,
	channel pkgamqp.AMQP,
) AMQPPublisher {
	return &AMQPPublisherImpl{
		logger:  logger,
		tracing: tracing,
		channel: channel,
	}
}

// Create - publish amqp create
func (publisherImpl *AMQPPublisherImpl) Create(value string) {
	ctx := context.Background()

	messageName := "send_email"

	// Create a new span (child of the trace id) to inform the publishing of the message
	tr := publisherImpl.tracing.Tracer("amqp")
	spanName := fmt.Sprintf("AMQP - publish - %s", messageName)

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
	}

	ctx, span := tr.Start(ctx, spanName, opts...)
	defer span.End()

	channel := publisherImpl.channel.Get()
	q, err := channel.QueueDeclare(messageName, true, false, false, false, nil)
	if err != nil {
		publisherImpl.logger.Error(err)
	}

	// Inject the context in the headers
	headers := pkgamqp.InjectAMQPHeaders(ctx)
	body := value
	msg := amqp.Publishing{
		Headers:     headers,
		ContentType: "text/plain",
		Body:        []byte(body),
	}

	err = channel.Publish("", q.Name, false, false, msg)
	if err != nil {
		publisherImpl.logger.Error(err)
	}
	publisherImpl.logger.Println("Publisher send to queue name", messageName)
}
