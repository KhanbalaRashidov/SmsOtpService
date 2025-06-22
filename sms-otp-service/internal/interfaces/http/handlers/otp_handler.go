package handlers

import (
	"sms-otp-service/internal/application/dto"
	"sms-otp-service/internal/application/usecases"
	"sms-otp-service/internal/domain/entities"
	"sms-otp-service/internal/domain/services"
	"sms-otp-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OTPHandler struct {
	otpUseCase     usecases.OTPUseCase
	phoneValidator *utils.PhoneValidator
	logger         *logrus.Logger
}

func NewOTPHandler(otpUseCase usecases.OTPUseCase, logger *logrus.Logger) *OTPHandler {
	return &OTPHandler{
		otpUseCase:     otpUseCase,
		phoneValidator: utils.NewPhoneValidator(),
		logger:         logger,
	}
}

// SendOTP godoc
// @Summary Send OTP
// @Description Send OTP to phone number
// @Tags OTP
// @Accept json
// @Produce json
// @Param request body dto.SendOTPRequest true "Send OTP request"
// @Success 200 {object} dto.SendOTPResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 429 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/otp/send [post]
func (h *OTPHandler) SendOTP(c *fiber.Ctx) error {
	var req dto.SendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "Invalid request format",
			Code:    "INVALID_REQUEST",
		})
	}

	if err := h.phoneValidator.Validate(req.PhoneNumber); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "Invalid phone number format",
			Code:    "INVALID_PHONE",
		})
	}

	req.PhoneNumber = h.phoneValidator.NormalizePhoneNumber(req.PhoneNumber)

	if req.Purpose == "" {
		req.Purpose = entities.PurposeVerification
	}

	resp, err := h.otpUseCase.SendOTP(c.Context(), &req)
	if err != nil {
		statusCode, errorResp := h.handleError(err)
		return c.Status(statusCode).JSON(errorResp)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verify OTP code
// @Tags OTP
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPRequest true "Verify OTP request"
// @Success 200 {object} dto.VerifyOTPResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/otp/verify [post]
func (h *OTPHandler) VerifyOTP(c *fiber.Ctx) error {
	var req dto.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "Invalid request format",
			Code:    "INVALID_REQUEST",
		})
	}

	if err := h.phoneValidator.Validate(req.PhoneNumber); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "Invalid phone number format",
			Code:    "INVALID_PHONE",
		})
	}

	if len(req.Code) != 6 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "OTP code must be 6 digits",
			Code:    "INVALID_CODE",
		})
	}

	req.PhoneNumber = h.phoneValidator.NormalizePhoneNumber(req.PhoneNumber)

	if req.Purpose == "" {
		req.Purpose = entities.PurposeVerification
	}

	resp, err := h.otpUseCase.VerifyOTP(c.Context(), &req)
	if err != nil {
		statusCode, errorResp := h.handleError(err)
		return c.Status(statusCode).JSON(errorResp)
	}

	if !resp.Success {
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// ResendOTP godoc
// @Summary Resend OTP
// @Description Resend OTP to phone number
// @Tags OTP
// @Accept json
// @Produce json
// @Param request body dto.ResendOTPRequest true "Resend OTP request"
// @Success 200 {object} dto.ResendOTPResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 429 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/otp/resend [post]
func (h *OTPHandler) ResendOTP(c *fiber.Ctx) error {
	var req dto.ResendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "Invalid request format",
			Code:    "INVALID_REQUEST",
		})
	}

	if err := h.phoneValidator.Validate(req.PhoneNumber); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Error:   "Invalid phone number format",
			Code:    "INVALID_PHONE",
		})
	}

	req.PhoneNumber = h.phoneValidator.NormalizePhoneNumber(req.PhoneNumber)

	if req.Purpose == "" {
		req.Purpose = entities.PurposeVerification
	}

	resp, err := h.otpUseCase.ResendOTP(c.Context(), &req)
	if err != nil {
		statusCode, errorResp := h.handleError(err)
		return c.Status(statusCode).JSON(errorResp)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *OTPHandler) handleError(err error) (int, dto.ErrorResponse) {
	switch err {
	case services.ErrRateLimitExceeded:
		return fiber.StatusTooManyRequests, dto.ErrorResponse{
			Success: false,
			Error:   "Rate limit exceeded. Please wait before requesting a new OTP.",
			Code:    "RATE_LIMIT",
		}
	case entities.ErrInvalidPhoneNumber:
		return fiber.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid phone number format",
			Code:    "INVALID_PHONE",
		}
	default:
		h.logger.WithError(err).Error("Unexpected error occurred")
		return fiber.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Internal server error",
			Code:    "INTERNAL_ERROR",
		}
	}
}
