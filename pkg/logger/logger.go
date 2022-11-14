package pkg_logger

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Error(err error)
	Println(args ...interface{})
	Printf(format string, args ...interface{})
}

type LoggerImpl struct{}

func NewLogger() Logger {
	return &LoggerImpl{}
}

func (logger *LoggerImpl) Error(err error) {
	logrus.Error(err)
}

func (logger *LoggerImpl) Println(args ...interface{}) {
	logrus.Println(args...)
}

func (logger *LoggerImpl) Printf(format string, args ...interface{}) {
	logrus.Printf(format, args...)
}
