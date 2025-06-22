package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"sms-otp-service/internal/domain/entities"
	"sms-otp-service/internal/domain/repositories"
	"time"
)

type gormOTPRepository struct {
	db *gorm.DB
}

func NewGormOTPRepository(db *gorm.DB) repositories.OTPRepository {
	return &gormOTPRepository{db: db}
}

func (r *gormOTPRepository) Create(ctx context.Context, otp *entities.OTP) error {
	return r.db.WithContext(ctx).Create(otp).Error
}

func (r *gormOTPRepository) FindByPhoneAndPurpose(ctx context.Context, phoneNumber string, purpose entities.OTPPurpose) (*entities.OTP, error) {
	var otp entities.OTP
	err := r.db.WithContext(ctx).
		Where("phone_number = ? AND purpose = ?", phoneNumber, purpose).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrOTPNotFound
		}
		return nil, err
	}

	return &otp, nil
}

func (r *gormOTPRepository) FindByID(ctx context.Context, id string) (*entities.OTP, error) {
	var otp entities.OTP
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&otp).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrOTPNotFound
		}
		return nil, err
	}

	return &otp, nil
}

func (r *gormOTPRepository) Update(ctx context.Context, otp *entities.OTP) error {
	otp.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(otp).Error
}

func (r *gormOTPRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entities.OTP{}, "id = ?", id).Error
}

func (r *gormOTPRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entities.OTP{}).Error
}

func (r *gormOTPRepository) FindActiveByPhone(ctx context.Context, phoneNumber string) ([]*entities.OTP, error) {
	var otps []*entities.OTP
	err := r.db.WithContext(ctx).
		Where("phone_number = ? AND expires_at > ? AND is_verified = false AND attempts < max_attempts",
			phoneNumber, time.Now()).
		Order("created_at DESC").
		Find(&otps).Error

	return otps, err
}

func (r *gormOTPRepository) CountRecentOTPs(ctx context.Context, phoneNumber string, minutes int) (int64, error) {
	var count int64
	cutoffTime := time.Now().Add(-time.Duration(minutes) * time.Minute)

	err := r.db.WithContext(ctx).
		Model(&entities.OTP{}).
		Where("phone_number = ? AND created_at > ?", phoneNumber, cutoffTime).
		Count(&count).Error

	return count, err
}
