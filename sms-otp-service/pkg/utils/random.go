package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

type OTPGenerator struct {
	length int
}

func NewOTPGenerator(length int) *OTPGenerator {
	return &OTPGenerator{length: length}
}

func (g *OTPGenerator) Generate() string {
	digits := make([]byte, g.length)
	for i := 0; i < g.length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		digits[i] = byte(num.Int64()) + '0'
	}
	return string(digits)
}

type PhoneValidator struct {
	patterns []*regexp.Regexp
}

func NewPhoneValidator() *PhoneValidator {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`^\+\d{10,15}$`),
		regexp.MustCompile(`^\+994(50|51|55|70|77|99)\d{7}$`),
		regexp.MustCompile(`^0(50|51|55|70|77|99)\d{7}$`),
	}

	return &PhoneValidator{patterns: patterns}
}

func (v *PhoneValidator) Validate(phoneNumber string) error {
	cleaned := v.cleanPhoneNumber(phoneNumber)

	if len(cleaned) < 10 || len(cleaned) > 15 {
		return fmt.Errorf("phone number length must be between 10 and 15 digits")
	}

	for _, pattern := range v.patterns {
		if pattern.MatchString(cleaned) {
			return nil
		}
	}

	return fmt.Errorf("invalid phone number format")
}

func (v *PhoneValidator) cleanPhoneNumber(phoneNumber string) string {
	cleaned := regexp.MustCompile(`[\s\-\(\)]`).ReplaceAllString(phoneNumber, "")

	if strings.HasPrefix(cleaned, "0") && len(cleaned) == 10 {
		cleaned = "+994" + cleaned[1:]
	}

	return cleaned
}

func (v *PhoneValidator) NormalizePhoneNumber(phoneNumber string) string {
	cleaned := v.cleanPhoneNumber(phoneNumber)

	if !strings.HasPrefix(cleaned, "+") {
		if strings.HasPrefix(cleaned, "994") {
			cleaned = "+" + cleaned
		} else if strings.HasPrefix(cleaned, "0") && len(cleaned) == 10 {
			cleaned = "+994" + cleaned[1:]
		}
	}

	return cleaned
}
