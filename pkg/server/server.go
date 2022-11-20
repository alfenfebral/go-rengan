package server

import (
	"context"
	pkg_amqp "go-rengan/pkg/amqp"
	pkg_logger "go-rengan/pkg/logger"
	pkg_mongodb "go-rengan/pkg/mongodb"
	pkg_http_server "go-rengan/pkg/server/http"
	pkg_tracing "go-rengan/pkg/tracing"
	todo_amqp_delivery "go-rengan/todo/delivery/amqp"

	"github.com/sirupsen/logrus"
)

type ServerImpl struct {
	httpServer       pkg_http_server.HTTPServer
	logger           pkg_logger.Logger
	Tracing          pkg_tracing.Tracing
	TodoAMQPConsumer todo_amqp_delivery.TodoAMQPConsumer
	MongoDB          pkg_mongodb.MongoDB
	AMQP             pkg_amqp.AMQP
}

func NewServer(
	tracing pkg_tracing.Tracing,
	logger pkg_logger.Logger,
	amqp pkg_amqp.AMQP,
	todoAMQPConsumer todo_amqp_delivery.TodoAMQPConsumer,
	mongoDB pkg_mongodb.MongoDB,
	httpServer pkg_http_server.HTTPServer,
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
func (serverImpl *ServerImpl) Run() error {
	go func() {
		serverImpl.TodoAMQPConsumer.Register()
	}()

	go func() {
		err := serverImpl.httpServer.Run()
		if err != nil {
			serverImpl.logger.Error(err)
		}
	}()

	return nil
}

// GracefulStop server
func (serverImpl *ServerImpl) GracefulStop(ctx context.Context, done chan bool) {
	err := serverImpl.httpServer.GracefulStop(ctx)
	if err != nil {
		serverImpl.logger.Error(err)
	}

	logrus.Info("gracefully shutdowned")
	done <- true
}
