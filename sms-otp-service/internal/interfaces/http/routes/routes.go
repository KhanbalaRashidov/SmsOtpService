package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"sms-otp-service/internal/interfaces/http/handlers"
)

type Routes struct {
	otpHandler    *handlers.OTPHandler
	healthHandler *handlers.HealthHandler
}

func NewRoutes(otpHandler *handlers.OTPHandler, healthHandler *handlers.HealthHandler) *Routes {
	return &Routes{
		otpHandler:    otpHandler,
		healthHandler: healthHandler,
	}
}

func (r *Routes) Setup(app *fiber.App) {
	// Middleware
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Get("/health", r.healthHandler.Health)
	app.Get("/ready", r.healthHandler.Ready)

	v1 := app.Group("/api/v1")

	otp := v1.Group("/otp")
	otp.Post("/send", r.otpHandler.SendOTP)
	otp.Post("/verify", r.otpHandler.VerifyOTP)
	otp.Post("/resend", r.otpHandler.ResendOTP)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "SMS OTP Service",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Route not found",
			"code":    "ROUTE_NOT_FOUND",
		})
	})
}
