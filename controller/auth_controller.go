package controller

import (
	"net/http"

	"auth-service/model"
	"auth-service/utils"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/auth/token", IssueTokens)
	app.Post("/auth/refresh", RefreshTokens)
}

func IssueTokens(c *fiber.Ctx) error {
	var request struct {
		UserID string `json:"user_id" validate:"required,uuid"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Validate request
	validate := utils.GetValidator()
	if err := validate.Struct(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Issue tokens
	accessToken, refreshToken, err := utils.GenerateTokens(request.UserID, c.IP())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate tokens"})
	}

	// Save refresh token hash
	if err := model.SaveRefreshToken(request.UserID, refreshToken, c.IP()); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not save refresh token"})
	}

	return c.JSON(fiber.Map{"access_token": accessToken, "refresh_token": refreshToken})
}

func RefreshTokens(c *fiber.Ctx) error {
	var request struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Validate token
	userID, ip, err := model.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		utils.EmailWarning("Invalid refresh token attempt") // log warning
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	// Check IP address
	if ip != c.IP() {
		utils.EmailWarning("IP mismatch detected during refresh") // log email warning
	}

	// Generate new tokens
	accessToken, refreshToken, err := utils.GenerateTokens(userID, c.IP())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate tokens"})
	}

	// Update refresh token hash
	if err := model.SaveRefreshToken(userID, refreshToken, c.IP()); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not save refresh token"})
	}

	return c.JSON(fiber.Map{"access_token": accessToken, "refresh_token": refreshToken})
}
