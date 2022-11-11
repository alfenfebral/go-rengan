package pkg_logger

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Error(err error)
}

type LoggerImpl struct{}

func NewLogger() Logger {
	return &LoggerImpl{}
}

func (logger *LoggerImpl) Error(err error) {
	logrus.Error(err)
}
