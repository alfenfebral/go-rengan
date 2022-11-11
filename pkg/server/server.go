package server

import (
	"context"
	pkg_logger "go-rengan/pkg/logger"
	pkg_http_server "go-rengan/pkg/server/http"

	"github.com/sirupsen/logrus"
)

type Server interface {
	Run() error
	GracefulStop(ctx context.Context, done chan bool)
}

type ServerImpl struct {
	httpServer pkg_http_server.HTTPServer
	logger     pkg_logger.Logger
}

func NewServer(logger pkg_logger.Logger, httpServer pkg_http_server.HTTPServer) Server {
	return &ServerImpl{
		httpServer: httpServer,
		logger:     logger,
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
