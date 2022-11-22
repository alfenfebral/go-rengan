package server

import (
	"context"
	amqp "go-rengan/pkg/amqp"
	logger "go-rengan/pkg/logger"
	mongodb "go-rengan/pkg/mongodb"
	httpserver "go-rengan/pkg/server/http"
	tracing "go-rengan/pkg/tracing"
	todoamqpdelivery "go-rengan/todo/delivery/amqp"

	"github.com/sirupsen/logrus"
)

type ServerImpl struct {
	httpServer       httpserver.HTTPServer
	logger           logger.Logger
	Tracing          tracing.Tracing
	TodoAMQPConsumer todoamqpdelivery.AMQPConsumer
	MongoDB          mongodb.MongoDB
	AMQP             amqp.AMQP
}

func NewServer(
	tracing tracing.Tracing,
	logger logger.Logger,
	amqp amqp.AMQP,
	todoAMQPConsumer todoamqpdelivery.AMQPConsumer,
	mongoDB mongodb.MongoDB,
	httpServer httpserver.HTTPServer,
) *ServerImpl {
	return &ServerImpl{
		httpServer:       httpServer,
		logger:           logger,
		Tracing:          tracing,
		AMQP:             amqp,
		TodoAMQPConsumer: todoAMQPConsumer,
		MongoDB:          mongoDB,
	}
}

// Run server
func (s *ServerImpl) Run() error {
	go func() {
		s.TodoAMQPConsumer.Register()
	}()

	go func() {
		err := s.httpServer.Run()
		if err != nil {
			s.logger.Error(err)
		}
	}()

	return nil
}

// GracefulStop server
func (s *ServerImpl) GracefulStop(ctx context.Context, done chan bool) {
	err := s.httpServer.GracefulStop(ctx)
	if err != nil {
		s.logger.Error(err)
	}

	logrus.Info("gracefully shutdowned")
	done <- true
}
