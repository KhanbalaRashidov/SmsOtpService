package usecases

import (
	"context"
	"fmt"
	"sms-otp-service/internal/application/dto"
	"sms-otp-service/internal/domain/entities"
	"sms-otp-service/internal/domain/services"
	"time"

	"github.com/sirupsen/logrus"
)

type OTPUseCase interface {
	SendOTP(ctx context.Context, req *dto.SendOTPRequest) (*dto.SendOTPResponse, error)
	VerifyOTP(ctx context.Context, req *dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error)
	ResendOTP(ctx context.Context, req *dto.ResendOTPRequest) (*dto.ResendOTPResponse, error)
}

type otpUseCase struct {
	otpDomainService services.OTPDomainService
	smsService       SMSService
	validityMinutes  int
	logger           *logrus.Logger
}

type SMSService interface {
	SendSMS(ctx context.Context, phoneNumber, message string) error
}

func NewOTPUseCase(
	otpDomainService services.OTPDomainService,
	smsService SMSService,
	validityMinutes int,
	logger *logrus.Logger,
) OTPUseCase {
	return &otpUseCase{
		otpDomainService: otpDomainService,
		smsService:       smsService,
		validityMinutes:  validityMinutes,
		logger:           logger,
	}
}

func (uc *otpUseCase) SendOTP(ctx context.Context, req *dto.SendOTPRequest) (*dto.SendOTPResponse, error) {
	if req.Purpose == "" {
		req.Purpose = entities.PurposeVerification
	}

	uc.logger.WithFields(logrus.Fields{
		"phone_number": req.PhoneNumber,
		"purpose":      req.Purpose,
	}).Info("Generating OTP")

	otp, err := uc.otpDomainService.GenerateOTP(ctx, req.PhoneNumber, req.Purpose)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to generate OTP")
		return nil, err
	}

	message := uc.buildSMSMessage(otp.Code, req.Purpose)
	if err := uc.smsService.SendSMS(ctx, req.PhoneNumber, message); err != nil {
		uc.logger.WithError(err).Error("Failed to send SMS")
		return nil, fmt.Errorf("failed to send SMS: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"otp_id":       otp.ID,
		"phone_number": req.PhoneNumber,
	}).Info("OTP sent successfully")

	return &dto.SendOTPResponse{
		Success:   true,
		Message:   "OTP sent successfully",
		ExpiresIn: uc.validityMinutes * 60,
		ID:        otp.ID.String(),
	}, nil
}

func (uc *otpUseCase) VerifyOTP(ctx context.Context, req *dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error) {
	if req.Purpose == "" {
		req.Purpose = entities.PurposeVerification
	}

	uc.logger.WithFields(logrus.Fields{
		"phone_number": req.PhoneNumber,
		"purpose":      req.Purpose,
	}).Info("Verifying OTP")

	err := uc.otpDomainService.VerifyOTP(ctx, req.PhoneNumber, req.Code, req.Purpose)
	if err != nil {
		uc.logger.WithError(err).Warn("OTP verification failed")
		return &dto.VerifyOTPResponse{
			Success: false,
			Message: uc.getErrorMessage(err),
		}, nil
	}

	uc.logger.WithFields(logrus.Fields{
		"phone_number": req.PhoneNumber,
		"purpose":      req.Purpose,
	}).Info("OTP verified successfully")

	return &dto.VerifyOTPResponse{
		Success:    true,
		Message:    "OTP verified successfully",
		VerifiedAt: time.Now(),
	}, nil
}

func (uc *otpUseCase) ResendOTP(ctx context.Context, req *dto.ResendOTPRequest) (*dto.ResendOTPResponse, error) {
	if req.Purpose == "" {
		req.Purpose = entities.PurposeVerification
	}

	uc.logger.WithFields(logrus.Fields{
		"phone_number": req.PhoneNumber,
		"purpose":      req.Purpose,
	}).Info("Resending OTP")

	otp, err := uc.otpDomainService.ResendOTP(ctx, req.PhoneNumber, req.Purpose)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to resend OTP")
		return nil, err
	}

	message := uc.buildSMSMessage(otp.Code, req.Purpose)
	if err := uc.smsService.SendSMS(ctx, req.PhoneNumber, message); err != nil {
		uc.logger.WithError(err).Error("Failed to send SMS")
		return nil, fmt.Errorf("failed to send SMS: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"otp_id":       otp.ID,
		"phone_number": req.PhoneNumber,
	}).Info("OTP resent successfully")

	return &dto.ResendOTPResponse{
		Success:   true,
		Message:   "OTP resent successfully",
		ExpiresIn: uc.validityMinutes * 60,
	}, nil
}

func (uc *otpUseCase) buildSMSMessage(code string, purpose entities.OTPPurpose) string {
	switch purpose {
	case entities.PurposeLogin:
		return fmt.Sprintf("Your login code is: %s. Valid for %d minutes. Do not share this code.", code, uc.validityMinutes)
	case entities.PurposeReset:
		return fmt.Sprintf("Your password reset code is: %s. Valid for %d minutes. Do not share this code.", code, uc.validityMinutes)
	default:
		return fmt.Sprintf("Your verification code is: %s. Valid for %d minutes. Do not share this code.", code, uc.validityMinutes)
	}
}

func (uc *otpUseCase) getErrorMessage(err error) string {
	switch err {
	case entities.ErrOTPExpired:
		return "OTP has expired. Please request a new one."
	case entities.ErrOTPAlreadyUsed:
		return "OTP has already been used. Please request a new one."
	case entities.ErrMaxAttemptsReached:
		return "Maximum verification attempts reached. Please request a new OTP."
	case entities.ErrInvalidOTPCode:
		return "Invalid OTP code. Please try again."
	case entities.ErrOTPNotFound:
		return "OTP not found. Please request a new one."
	case services.ErrRateLimitExceeded:
		return "Too many requests. Please wait before requesting a new OTP."
	default:
		return "Verification failed. Please try again."
	}
}
