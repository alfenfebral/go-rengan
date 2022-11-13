//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package dep

import (
	pkg_amqp "go-rengan/pkg/amqp"
	pkg_logger "go-rengan/pkg/logger"
	pkg_mongodb "go-rengan/pkg/mongodb"
	pkg_server "go-rengan/pkg/server"
	pkg_http_server "go-rengan/pkg/server/http"
	pkg_tracing "go-rengan/pkg/tracing"
	todo_http "go-rengan/todo/delivery/http"
	repository "go-rengan/todo/repository"
	service "go-rengan/todo/service"

	"github.com/google/wire"
)

func InitializeServer() (*pkg_server.ServerImpl, error) {
	wire.Build(
		pkg_amqp.NewAMQP,
		pkg_tracing.NewTracing,
		pkg_logger.NewLogger,
		pkg_mongodb.NewMongoDB,
		repository.NewMongoTodoRepository,
		service.NewTodoService,
		pkg_http_server.NewHTTPServer,
		todo_http.NewTodoHTTPHandler,
		pkg_server.NewServer,
	)

	return &pkg_server.ServerImpl{}, nil
}
