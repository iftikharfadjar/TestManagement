package rest

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"boilerplate/services/test/domain"
)

type TestHandler struct {
	uc domain.TestUseCase
}

func NewTestHandler(uc domain.TestUseCase) *TestHandler {
	return &TestHandler{uc: uc}
}

func (h *TestHandler) SetupRoutes(router fiber.Router) {
	v1 := router.Group("/api/v1")
	v1.Get("/test", h.List)
	v1.Post("/test", h.Create)
	v1.Get("/test/:id", h.GetByID)
	v1.Put("/test/:id", h.Update)
	v1.Delete("/test/:id", h.Delete)
}

type TestRequest struct {
	TestID      string `json:"test_id"`
	TestName    string `json:"test_name"`
	Description string `json:"description"`
	Steps       string `json:"steps"`
	Status      string `json:"status"`
	Remarks     string `json:"remarks"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedBy   string `json:"updated_by"`
	UpdatedAt   string `json:"updated_at"`
	IsActive    bool   `json:"is_active"`
}

func (h *TestHandler) Create(c fiber.Ctx) error {
	var req TestRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	now := time.Now()
	test := &domain.Test{
		TestID:      req.TestID,
		TestName:    req.TestName,
		Description: req.Description,
		Steps:       req.Steps,
		Status:      req.Status,
		Remarks:     req.Remarks,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   now,
		UpdatedBy:   req.UpdatedBy,
		UpdatedAt:   now,
		IsActive:    req.IsActive,
	}

	created, err := h.uc.Create(c.Context(), test)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *TestHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")

	test, err := h.uc.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "test not found"})
	}

	return c.JSON(test)
}

func (h *TestHandler) List(c fiber.Ctx) error {
	tests, err := h.uc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(tests)
}

func (h *TestHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")

	var req TestRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	test := &domain.Test{
		TestName:    req.TestName,
		Description: req.Description,
		Steps:       req.Steps,
		Status:      req.Status,
		Remarks:     req.Remarks,
		UpdatedBy:   req.UpdatedBy,
		UpdatedAt:   time.Now(),
		IsActive:    req.IsActive,
	}

	updated, err := h.uc.Update(c.Context(), id, test)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updated)
}

func (h *TestHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	if err := h.uc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "test deleted successfully"})
}
