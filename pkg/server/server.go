package server

import (
	"context"
	pkg_amqp "go-rengan/pkg/amqp"
	pkg_logger "go-rengan/pkg/logger"
	pkg_mongodb "go-rengan/pkg/mongodb"
	pkg_http_server "go-rengan/pkg/server/http"
	pkg_tracing "go-rengan/pkg/tracing"
	todo_amqp "go-rengan/todo/delivery/amqp"

	"github.com/sirupsen/logrus"
)

type ServerImpl struct {
	httpServer       pkg_http_server.HTTPServer
	logger           pkg_logger.Logger
	Tp               pkg_tracing.Tracing
	TodoAMQPConsumer todo_amqp.TodoAMQPConsumer
	MongoDB          pkg_mongodb.MongoDB
	AMQP             pkg_amqp.AMQP
}

func NewServer(
	tp pkg_tracing.Tracing,
	logger pkg_logger.Logger,
	amqp pkg_amqp.AMQP,
	todoAMQP todo_amqp.TodoAMQPConsumer,
	mongoDB pkg_mongodb.MongoDB,
	httpServer pkg_http_server.HTTPServer,
) *ServerImpl {
	return &ServerImpl{
		httpServer:       httpServer,
		logger:           logger,
		Tp:               tp,
		AMQP:             amqp,
		TodoAMQPConsumer: todoAMQP,
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
