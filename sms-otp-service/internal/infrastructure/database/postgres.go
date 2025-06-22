package database

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sms-otp-service/internal/domain/entities"
	"sms-otp-service/internal/infrastructure/config"
	"time"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	var gormLogger logger.Interface
	if cfg.Logger.Level == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.Info("Database connection established successfully")

	return &Database{DB: db}, nil
}

func (d *Database) AutoMigrate() error {
	logrus.Info("Starting database migration...")

	if err := d.DB.AutoMigrate(&entities.OTP{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	if err := d.createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	logrus.Info("Database migration completed successfully")
	return nil
}

func (d *Database) createIndexes() error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_otps_phone_purpose ON otps(phone_number, purpose)",
		"CREATE INDEX IF NOT EXISTS idx_otps_phone_expires ON otps(phone_number, expires_at)",
		"CREATE INDEX IF NOT EXISTS idx_otps_created_at ON otps(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_otps_verified ON otps(is_verified, expires_at)",
	}

	for _, index := range indexes {
		if err := d.DB.Exec(index).Error; err != nil {
			logrus.Warnf("Failed to create index: %s, error: %v", index, err)
		}
	}

	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) Health() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}
