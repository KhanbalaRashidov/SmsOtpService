package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	SMS      SMSConfig
	OTP      OTPConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	DSN      string
}

type SMSConfig struct {
	Provider    string
	APIKey      string
	APISecret   string
	SenderName  string
	APIEndpoint string
}

type OTPConfig struct {
	ValidityMinutes  int
	RateLimitMinutes int
	MaxOTPsPerPeriod int
	CodeLength       int
	MaxAttempts      int
	CleanupInterval  time.Duration
}

type LoggerConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	cfg := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  parseDuration(getEnv("SERVER_READ_TIMEOUT", "30s")),
			WriteTimeout: parseDuration(getEnv("SERVER_WRITE_TIMEOUT", "30s")),
			IdleTimeout:  parseDuration(getEnv("SERVER_IDLE_TIMEOUT", "120s")),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "sms_otp_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		SMS: SMSConfig{
			Provider:    getEnv("SMS_PROVIDER", "mock"),
			APIKey:      getEnv("SMS_API_KEY", ""),
			APISecret:   getEnv("SMS_API_SECRET", ""),
			SenderName:  getEnv("SMS_SENDER_NAME", "OTPService"),
			APIEndpoint: getEnv("SMS_API_ENDPOINT", ""),
		},
		OTP: OTPConfig{
			ValidityMinutes:  parseInt(getEnv("OTP_VALIDITY_MINUTES", "5")),
			RateLimitMinutes: parseInt(getEnv("OTP_RATE_LIMIT_MINUTES", "10")),
			MaxOTPsPerPeriod: parseInt(getEnv("OTP_MAX_PER_PERIOD", "3")),
			CodeLength:       parseInt(getEnv("OTP_CODE_LENGTH", "6")),
			MaxAttempts:      parseInt(getEnv("OTP_MAX_ATTEMPTS", "3")),
			CleanupInterval:  parseDuration(getEnv("OTP_CLEANUP_INTERVAL", "1h")),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	cfg.Database.DSN = buildDSN(cfg.Database)

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}

func buildDSN(cfg DatabaseConfig) string {
	return "host=" + cfg.Host +
		" port=" + cfg.Port +
		" user=" + cfg.User +
		" password=" + cfg.Password +
		" dbname=" + cfg.DBName +
		" sslmode=" + cfg.SSLMode
}
