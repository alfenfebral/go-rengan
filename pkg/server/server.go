package server

import (
	"context"
	pkg_logger "go-rengan/pkg/logger"
	pkg_http_server "go-rengan/pkg/server/http"
	pkg_tracing "go-rengan/pkg/tracing"

	"github.com/sirupsen/logrus"
)

type ServerImpl struct {
	httpServer pkg_http_server.HTTPServer
	logger     pkg_logger.Logger
	Tp         pkg_tracing.Tracing
}

func NewServer(
	tp pkg_tracing.Tracing,
	logger pkg_logger.Logger,
	httpServer pkg_http_server.HTTPServer,
) *ServerImpl {
	return &ServerImpl{
		httpServer: httpServer,
		logger:     logger,
		Tp:         tp,
	}
}

// Run server
func (s *ServerImpl) Run() error {
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
