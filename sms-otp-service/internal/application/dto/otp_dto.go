package dto

import (
	"sms-otp-service/internal/domain/entities"
	"time"
)

type SendOTPRequest struct {
	PhoneNumber string              `json:"phone_number" validate:"required,phone"`
	Purpose     entities.OTPPurpose `json:"purpose,omitempty"`
}

type SendOTPResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`
	ID        string `json:"id,omitempty"`
}

type VerifyOTPRequest struct {
	PhoneNumber string              `json:"phone_number" validate:"required,phone"`
	Code        string              `json:"code" validate:"required,len=6,numeric"`
	Purpose     entities.OTPPurpose `json:"purpose,omitempty"`
}

type VerifyOTPResponse struct {
	Success    bool      `json:"success"`
	Message    string    `json:"message"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
}

type ResendOTPRequest struct {
	PhoneNumber string              `json:"phone_number" validate:"required,phone"`
	Purpose     entities.OTPPurpose `json:"purpose,omitempty"`
}

type ResendOTPResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}
