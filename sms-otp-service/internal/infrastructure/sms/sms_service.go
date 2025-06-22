package sms

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"sms-otp-service/internal/infrastructure/config"
)

type Service interface {
	SendSMS(ctx context.Context, phoneNumber, message string) error
}

func NewSMSService(cfg *config.Config, logger *logrus.Logger) Service {
	switch cfg.SMS.Provider {
	case "mock":
		return NewMockSMSService(logger)
	default:
		logger.Warn("Unknown SMS provider, falling back to mock")
		return NewMockSMSService(logger)
	}
}

type mockSMSService struct {
	logger *logrus.Logger
}

func NewMockSMSService(logger *logrus.Logger) Service {
	return &mockSMSService{
		logger: logger,
	}
}

func (s *mockSMSService) SendSMS(ctx context.Context, phoneNumber, message string) error {
	s.logger.WithFields(logrus.Fields{
		"phone_number": phoneNumber,
		"message":      message,
	}).Info("📱 Mock SMS sent")

	// Simulate SMS sending
	fmt.Printf("\n=== MOCK SMS ===\n")
	fmt.Printf("To: %s\n", phoneNumber)
	fmt.Printf("Message: %s\n", message)
	fmt.Printf("================\n\n")

	return nil
}
