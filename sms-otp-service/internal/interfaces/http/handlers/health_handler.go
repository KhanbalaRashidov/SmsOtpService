package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"sms-otp-service/internal/application/dto"
	"sms-otp-service/internal/infrastructure/database"
	"time"
)

type HealthHandler struct {
	db     *database.Database
	logger *logrus.Logger
}

func NewHealthHandler(db *database.Database, logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: logger,
	}
}

// Health godoc
// @Summary Health check
// @Description Get health status of the service
// @Tags Health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	services := make(map[string]string)
	overallStatus := "healthy"

	// Check database
	if err := h.db.Health(); err != nil {
		services["database"] = "unhealthy"
		overallStatus = "unhealthy"
		h.logger.WithError(err).Error("Database health check failed")
	} else {
		services["database"] = "healthy"
	}

	services["sms"] = "healthy" // Mock SMS is always healthy

	response := dto.HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Services:  services,
		Version:   "1.0.0",
	}

	if overallStatus == "unhealthy" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(response)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Ready godoc
// @Summary Readiness check
// @Description Check if service is ready to serve requests
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ready [get]
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	// Check if all required services are available
	if err := h.db.Health(); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "not ready",
			"error":  "database unavailable",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ready",
	})
}
