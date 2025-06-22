package entities

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

var (
	ErrOTPExpired         = errors.New("otp has expired")
	ErrOTPAlreadyUsed     = errors.New("otp has already been used")
	ErrMaxAttemptsReached = errors.New("maximum verification attempts reached")
	ErrInvalidOTPCode     = errors.New("invalid otp code")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrOTPNotFound        = errors.New("otp not found")
)

type OTPPurpose string

const (
	PurposeVerification OTPPurpose = "verification"
	PurposeLogin        OTPPurpose = "login"
	PurposeReset        OTPPurpose = "reset"
)

type OTP struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PhoneNumber string     `json:"phone_number" gorm:"type:varchar(20);not null;index"`
	Code        string     `json:"code" gorm:"type:varchar(6);not null;index"`
	Purpose     OTPPurpose `json:"purpose" gorm:"type:varchar(50);not null;default:'verification'"`
	IsVerified  bool       `json:"is_verified" gorm:"default:false"`
	Attempts    int        `json:"attempts" gorm:"default:0"`
	MaxAttempts int        `json:"max_attempts" gorm:"default:3"`
	ExpiresAt   time.Time  `json:"expires_at" gorm:"not null;index"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`
}

func (OTP) TableName() string {
	return "otps"
}

func (o *OTP) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func NewOTP(phoneNumber, code string, purpose OTPPurpose, validityMinutes int) *OTP {
	now := time.Now()
	return &OTP{
		ID:          uuid.New(),
		PhoneNumber: phoneNumber,
		Code:        code,
		Purpose:     purpose,
		IsVerified:  false,
		Attempts:    0,
		MaxAttempts: 3,
		ExpiresAt:   now.Add(time.Duration(validityMinutes) * time.Minute),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

func (o *OTP) CanAttempt() bool {
	return o.Attempts < o.MaxAttempts
}

func (o *OTP) Verify(code string) error {
	if o.IsVerified {
		return ErrOTPAlreadyUsed
	}
	if o.IsExpired() {
		return ErrOTPExpired
	}
	if !o.CanAttempt() {
		return ErrMaxAttemptsReached
	}

	o.Attempts++
	o.UpdatedAt = time.Now()

	if o.Code != code {
		return ErrInvalidOTPCode
	}

	o.IsVerified = true
	now := time.Now()
	o.VerifiedAt = &now
	return nil
}

func (o *OTP) IsValid() bool {
	return !o.IsExpired() && !o.IsVerified && o.CanAttempt()
}
