// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package dep

import (
	"go-rengan/pkg/amqp"
	"go-rengan/pkg/logger"
	"go-rengan/pkg/mongodb"
	"go-rengan/pkg/server"
	"go-rengan/pkg/server/http"
	"go-rengan/pkg/tracing"
	"go-rengan/todo/delivery/amqp"
	"go-rengan/todo/delivery/http"
	"go-rengan/todo/repository"
	"go-rengan/todo/service"
)

// Injectors from wire.go:

func InitializeServer() (*server.ServerImpl, error) {
	tracing, err := pkg_tracing.NewTracing()
	if err != nil {
		return nil, err
	}
	logger := pkg_logger.NewLogger()
	amqp, err := pkg_amqp.NewAMQP()
	if err != nil {
		return nil, err
	}
	todoAMQPConsumer := todo_amqp.NewTodoAMQPConsumer(logger, tracing, amqp)
	mongoDB, err := mongodb.NewMongoDB()
	if err != nil {
		return nil, err
	}
	mongoTodoRepository := repository.NewMongoTodoRepository(mongoDB)
	todoService := service.NewTodoService(tracing, mongoTodoRepository)
	todoHTTPHandler := todo_http.NewTodoHTTPHandler(tracing, todoService)
	httpServer := pkg_http_server.NewHTTPServer(logger, todoHTTPHandler)
	serverImpl := server.NewServer(tracing, logger, amqp, todoAMQPConsumer, mongoDB, httpServer)
	return serverImpl, nil
}
