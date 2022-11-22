package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Error(args ...interface{})
	Println(args ...interface{})
	Printf(format string, args ...interface{})
}

type LoggerImpl struct{}

func New() Logger {
	return &LoggerImpl{}
}

func (l *LoggerImpl) Error(args ...interface{}) {
	logrus.Error(args...)
}

func (l *LoggerImpl) Println(args ...interface{}) {
	logrus.Println(args...)
}

func (l *LoggerImpl) Printf(format string, args ...interface{}) {
	logrus.Printf(format, args...)
}
