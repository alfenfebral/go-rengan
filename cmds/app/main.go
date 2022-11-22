package main

import (
	"context"
	"go-rengan/dep"
	config "go-rengan/pkg/config"
	logger "go-rengan/pkg/logger"
	validator "go-rengan/pkg/validator"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := logger.New()

	// Config
	config.New()

	// Validator
	validator.New()

	// Server
	server, err := dep.InitializeServer()
	if err != nil {
		logger.Error(err)
	}

	defer server.Tracing.ShutDown()
	defer server.AMQP.Get().Close()
	defer func() {
		if err := server.MongoDB.Disconnect(); err != nil {
			logger.Error(err)
		}
	}()

	go func() {
		err := server.Run()
		if err != nil {
			logger.Error(err)
		}
	}()

	// catch shutdown
	done := make(chan bool, 1)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		// graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.GracefulStop(ctx, done)
	}()

	// wait for graceful shutdown
	<-done
}
