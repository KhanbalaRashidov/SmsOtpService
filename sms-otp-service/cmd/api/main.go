package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"log"
	"os"
	"os/signal"
	_ "sms-otp-service/docs" // swagger docs
	"sms-otp-service/internal/application/usecases"
	"sms-otp-service/internal/domain/repositories"
	"sms-otp-service/internal/domain/services"
	"sms-otp-service/internal/infrastructure/config"
	"sms-otp-service/internal/infrastructure/database"
	infraRepos "sms-otp-service/internal/infrastructure/repositories"
	"sms-otp-service/internal/infrastructure/sms"
	"sms-otp-service/internal/interfaces/http/handlers"
	"sms-otp-service/internal/interfaces/http/routes"
	"sms-otp-service/pkg/logger"
	"sms-otp-service/pkg/utils"
	"syscall"
	"time"
)

// @title SMS OTP Service API
// @version 1.0
// @description Professional SMS OTP service with Clean Architecture
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@sms-otp-service.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg)
	appLogger.Info("Starting SMS OTP Service...")

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		appLogger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Auto migrate database
	if err := db.AutoMigrate(); err != nil {
		appLogger.WithError(err).Fatal("Failed to migrate database")
	}

	otpRepo := infraRepos.NewGormOTPRepository(db.DB)

	otpGenerator := utils.NewOTPGenerator(cfg.OTP.CodeLength)
	phoneValidator := utils.NewPhoneValidator()

	smsService := sms.NewSMSService(cfg, appLogger)

	otpDomainService := services.NewOTPDomainService(
		otpRepo,
		otpGenerator,
		phoneValidator,
		cfg.OTP.RateLimitMinutes,
		cfg.OTP.MaxOTPsPerPeriod,
		cfg.OTP.ValidityMinutes,
	)

	otpUseCase := usecases.NewOTPUseCase(
		otpDomainService,
		smsService,
		cfg.OTP.ValidityMinutes,
		appLogger,
	)

	otpHandler := handlers.NewOTPHandler(otpUseCase, appLogger)
	healthHandler := handlers.NewHealthHandler(db, appLogger)

	routesHandler := routes.NewRoutes(otpHandler, healthHandler)

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			appLogger.WithError(err).Error("Request error")

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
				"code":    "REQUEST_ERROR",
			})
		},
	})

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	routesHandler.Setup(app)

	go startCleanupRoutine(otpRepo, cfg.OTP.CleanupInterval, appLogger)

	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	go func() {
		appLogger.WithField("address", serverAddr).Info("Starting HTTP server...")
		if err := app.Listen(serverAddr); err != nil {
			appLogger.WithError(err).Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		appLogger.WithError(err).Error("Server forced to shutdown")
	}

	appLogger.Info("Server exited")
}

func startCleanupRoutine(otpRepo repositories.OTPRepository, interval time.Duration, logger *logrus.Logger) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.WithField("interval", interval).Info("Starting OTP cleanup routine")

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			if err := otpRepo.DeleteExpired(ctx); err != nil {
				logger.WithError(err).Error("Failed to clean up expired OTPs")
			} else {
				logger.Debug("Expired OTPs cleaned up successfully")
			}

			cancel()
		}
	}
}
