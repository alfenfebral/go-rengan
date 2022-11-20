package main

import (
	"context"
	"go-rengan/dep"
	pkg_config "go-rengan/pkg/config"
	pkg_validator "go-rengan/pkg/validator"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// Config
	pkg_config.NewConfig()

	// Validator
	pkg_validator.NewValidator()

	// Server
	server, err := dep.InitializeServer()
	if err != nil {
		logrus.Error(err)
	}

	defer server.Tracing.ShutDown()
	defer server.AMQP.Get().Close()
	defer func() {
		if err := server.MongoDB.Disconnect(); err != nil {
			logrus.Error(err)
		}
	}()

	go func() {
		err := server.Run()
		if err != nil {
			logrus.Error(err)
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
