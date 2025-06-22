package services

import (
	"context"
	"errors"
	"sms-otp-service/internal/domain/entities"
	"sms-otp-service/internal/domain/repositories"
)

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrInvalidRequest    = errors.New("invalid request")
)

type OTPDomainService interface {
	GenerateOTP(ctx context.Context, phoneNumber string, purpose entities.OTPPurpose) (*entities.OTP, error)
	VerifyOTP(ctx context.Context, phoneNumber, code string, purpose entities.OTPPurpose) error
	ResendOTP(ctx context.Context, phoneNumber string, purpose entities.OTPPurpose) (*entities.OTP, error)
}

type otpDomainService struct {
	otpRepo          repositories.OTPRepository
	otpGenerator     OTPGenerator
	phoneValidator   PhoneValidator
	rateLimitMinutes int
	maxOTPsPerPeriod int
	validityMinutes  int
}

type OTPGenerator interface {
	Generate() string
}

type PhoneValidator interface {
	Validate(phoneNumber string) error
}

func NewOTPDomainService(
	otpRepo repositories.OTPRepository,
	otpGenerator OTPGenerator,
	phoneValidator PhoneValidator,
	rateLimitMinutes, maxOTPsPerPeriod, validityMinutes int,
) OTPDomainService {
	return &otpDomainService{
		otpRepo:          otpRepo,
		otpGenerator:     otpGenerator,
		phoneValidator:   phoneValidator,
		rateLimitMinutes: rateLimitMinutes,
		maxOTPsPerPeriod: maxOTPsPerPeriod,
		validityMinutes:  validityMinutes,
	}
}

func (s *otpDomainService) GenerateOTP(ctx context.Context, phoneNumber string, purpose entities.OTPPurpose) (*entities.OTP, error) {
	if err := s.phoneValidator.Validate(phoneNumber); err != nil {
		return nil, entities.ErrInvalidPhoneNumber
	}

	recentCount, err := s.otpRepo.CountRecentOTPs(ctx, phoneNumber, s.rateLimitMinutes)
	if err != nil {
		return nil, err
	}

	if recentCount >= int64(s.maxOTPsPerPeriod) {
		return nil, ErrRateLimitExceeded
	}

	existingOTPs, err := s.otpRepo.FindActiveByPhone(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}

	for _, existingOTP := range existingOTPs {
		if existingOTP.Purpose == purpose {
			existingOTP.Attempts = existingOTP.MaxAttempts
			if err := s.otpRepo.Update(ctx, existingOTP); err != nil {
				return nil, err
			}
		}
	}

	code := s.otpGenerator.Generate()
	otp := entities.NewOTP(phoneNumber, code, purpose, s.validityMinutes)

	if err := s.otpRepo.Create(ctx, otp); err != nil {
		return nil, err
	}

	return otp, nil
}

func (s *otpDomainService) VerifyOTP(ctx context.Context, phoneNumber, code string, purpose entities.OTPPurpose) error {
	otp, err := s.otpRepo.FindByPhoneAndPurpose(ctx, phoneNumber, purpose)
	if err != nil {
		return entities.ErrOTPNotFound
	}

	if err := otp.Verify(code); err != nil {
		if updateErr := s.otpRepo.Update(ctx, otp); updateErr != nil {
			return updateErr
		}
		return err
	}

	if err := s.otpRepo.Update(ctx, otp); err != nil {
		return err
	}

	return nil
}

func (s *otpDomainService) ResendOTP(ctx context.Context, phoneNumber string, purpose entities.OTPPurpose) (*entities.OTP, error) {
	existingOTP, err := s.otpRepo.FindByPhoneAndPurpose(ctx, phoneNumber, purpose)
	if err == nil && !existingOTP.IsExpired() {
		recentCount, err := s.otpRepo.CountRecentOTPs(ctx, phoneNumber, 1)
		if err != nil {
			return nil, err
		}
		if recentCount > 0 {
			return nil, ErrRateLimitExceeded
		}
	}

	return s.GenerateOTP(ctx, phoneNumber, purpose)
}
