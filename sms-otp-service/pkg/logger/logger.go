package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"sms-otp-service/internal/infrastructure/config"
)

func NewLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	if cfg.Logger.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}

	logger.SetOutput(os.Stdout)

	return logger
}
