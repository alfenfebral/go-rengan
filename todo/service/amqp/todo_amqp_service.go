package todo_amqp_service

import (
	"context"
	"fmt"
	pkg_logger "go-rengan/pkg/logger"
	pkg_tracing "go-rengan/pkg/tracing"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"

	pkg_amqp "go-rengan/pkg/amqp"
)

type TodoAMQPPublisher interface {
	Create(value string)
}

type todoAMQPPublisher struct {
	logger  pkg_logger.Logger
	tp      pkg_tracing.Tracing
	channel pkg_amqp.AMQP
}

func NewTodoAMQPPublisher(
	logger pkg_logger.Logger,
	tp pkg_tracing.Tracing,
	channel pkg_amqp.AMQP,
) TodoAMQPPublisher {
	return &todoAMQPPublisher{
		logger:  logger,
		tp:      tp,
		channel: channel,
	}
}

// Create - publish amqp create
func (p *todoAMQPPublisher) Create(value string) {
	ctx := context.Background()

	messageName := "send_email"

	// Create a new span (child of the trace id) to inform the publishing of the message
	tr := p.tp.Tracer("amqp")
	spanName := fmt.Sprintf("AMQP - publish - %s", messageName)

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
	}

	ctx, span := tr.Start(ctx, spanName, opts...)
	defer span.End()

	channel := p.channel.Get()
	q, err := channel.QueueDeclare(messageName, true, false, false, false, nil)
	if err != nil {
		p.logger.Error(err)
	}

	// Inject the context in the headers
	headers := pkg_amqp.InjectAMQPHeaders(ctx)
	body := value
	msg := amqp.Publishing{
		Headers:     headers,
		ContentType: "text/plain",
		Body:        []byte(body),
	}

	err = channel.Publish("", q.Name, false, false, msg)
	if err != nil {
		p.logger.Error(err)
	}
	p.logger.Println("Publisher send to queue name", messageName)
}
