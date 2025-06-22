package repositories

import (
	"context"
	"sms-otp-service/internal/domain/entities"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *entities.OTP) error
	FindByPhoneAndPurpose(ctx context.Context, phoneNumber string, purpose entities.OTPPurpose) (*entities.OTP, error)

	FindByID(ctx context.Context, id string) (*entities.OTP, error)

	Update(ctx context.Context, otp *entities.OTP) error

	Delete(ctx context.Context, id string) error

	DeleteExpired(ctx context.Context) error

	FindActiveByPhone(ctx context.Context, phoneNumber string) ([]*entities.OTP, error)

	CountRecentOTPs(ctx context.Context, phoneNumber string, minutes int) (int64, error)
}
