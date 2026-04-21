package rest

import (
	"github.com/gofiber/fiber/v3"
	"boilerplate/services/auth/domain"
)

type AuthHandler struct {
	uc domain.AuthUseCase
}

func NewAuthHandler(uc domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) SetupRoutes(router fiber.Router) {
	v1 := router.Group("/api/v1")
	v1.Post("/signup", h.Signup)
	v1.Post("/login", h.Login)
	v1.Post("/logout", h.Logout)
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Signup(c fiber.Ctx) error {
	var req SignupRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	user, err := h.uc.Signup(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Signup successful", "user": user})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	token, err := h.uc.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.JSON(fiber.Map{"message": "Login successful", "token": token})
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}