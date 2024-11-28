package utils

import (
	"github.com/sirupsen/logrus"
)

func NewLogger(level string) *logrus.Logger {
	logger := logrus.New()

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parsedLevel = logrus.InfoLevel
	}
	logger.SetLevel(parsedLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}
