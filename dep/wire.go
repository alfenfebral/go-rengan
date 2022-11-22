//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package dep

import (
	amqp "go-rengan/pkg/amqp"
	logger "go-rengan/pkg/logger"
	mongodb "go-rengan/pkg/mongodb"
	server "go-rengan/pkg/server"
	httpserver "go-rengan/pkg/server/http"
	tracing "go-rengan/pkg/tracing"
	todoamqpdelivery "go-rengan/todo/delivery/amqp"
	todohttpdelivery "go-rengan/todo/delivery/http"
	todoamqpservice "go-rengan/todo/publisher"
	repository "go-rengan/todo/repository"
	service "go-rengan/todo/service"

	"github.com/google/wire"
)

func InitializeServer() (*server.ServerImpl, error) {
	wire.Build(
		amqp.New,
		tracing.New,
		logger.New,
		mongodb.New,
		repository.New,
		service.New,
		httpserver.New,
		todohttpdelivery.New,
		todoamqpdelivery.New,
		todoamqpservice.New,
		server.NewServer,
	)

	return &server.ServerImpl{}, nil
}
