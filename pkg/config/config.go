package pkg_config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func NewConfig() error {
	err := godotenv.Load()
	if err != nil {
		logrus.Error(err)
	}

	return nil
}
