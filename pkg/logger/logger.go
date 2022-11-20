package pkg_logger

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Error(args ...interface{})
	Println(args ...interface{})
	Printf(format string, args ...interface{})
}

type LoggerImpl struct{}

func NewLogger() Logger {
	return &LoggerImpl{}
}

func (logger *LoggerImpl) Error(args ...interface{}) {
	logrus.Error(args...)
}

func (logger *LoggerImpl) Println(args ...interface{}) {
	logrus.Println(args...)
}

func (logger *LoggerImpl) Printf(format string, args ...interface{}) {
	logrus.Printf(format, args...)
}
