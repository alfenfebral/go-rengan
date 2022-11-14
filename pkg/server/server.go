package server

import (
	"context"
	pkg_amqp "go-rengan/pkg/amqp"
	pkg_logger "go-rengan/pkg/logger"
	pkg_mongodb "go-rengan/pkg/mongodb"
	pkg_http_server "go-rengan/pkg/server/http"
	pkg_tracing "go-rengan/pkg/tracing"
	todo_amqp_delivery "go-rengan/todo/delivery/amqp"
	todo_amqp_service "go-rengan/todo/service/amqp"

	"github.com/sirupsen/logrus"
)

type ServerImpl struct {
	httpServer        pkg_http_server.HTTPServer
	logger            pkg_logger.Logger
	Tp                pkg_tracing.Tracing
	TodoAMQPConsumer  todo_amqp_delivery.TodoAMQPConsumer
	TodoAMQPPublisher todo_amqp_service.TodoAMQPPublisher
	MongoDB           pkg_mongodb.MongoDB
	AMQP              pkg_amqp.AMQP
}

func NewServer(
	tp pkg_tracing.Tracing,
	logger pkg_logger.Logger,
	amqp pkg_amqp.AMQP,
	todoAMQPConsumer todo_amqp_delivery.TodoAMQPConsumer,
	todoAMQPPublisher todo_amqp_service.TodoAMQPPublisher,
	mongoDB pkg_mongodb.MongoDB,
	httpServer pkg_http_server.HTTPServer,
) *ServerImpl {
	return &ServerImpl{
		httpServer:        httpServer,
		logger:            logger,
		Tp:                tp,
		AMQP:              amqp,
		TodoAMQPConsumer:  todoAMQPConsumer,
		TodoAMQPPublisher: todoAMQPPublisher,
		MongoDB:           mongoDB,
	}
}

// Run server
func (s *ServerImpl) Run() error {
	go func() {
		s.TodoAMQPPublisher.Create()
	}()

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
