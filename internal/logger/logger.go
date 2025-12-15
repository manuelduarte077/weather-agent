package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	// Set JSON formatter for structured logging (production-ready)
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Set log level from environment variable, default to info
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		Logger.SetLevel(logrus.InfoLevel)
		Logger.WithError(err).Warn("Invalid LOG_LEVEL, defaulting to info")
	} else {
		Logger.SetLevel(level)
	}

	// Set output to stdout
	Logger.SetOutput(os.Stdout)
}
