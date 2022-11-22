package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func New() error {
	err := godotenv.Load()
	if err != nil {
		logrus.Error(err)
	}

	return nil
}
