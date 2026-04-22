package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"

	"boilerplate/services/test/delivery/rest"
	"boilerplate/services/test/domain"
)

// --- Mock ---

type mockTestUseCase struct {
	createFn  func(ctx context.Context, test *domain.Test) (*domain.Test, error)
	getByIDFn func(ctx context.Context, id string) (*domain.Test, error)
	listFn    func(ctx context.Context) ([]domain.Test, error)
	updateFn  func(ctx context.Context, id string, test *domain.Test) (*domain.Test, error)
	deleteFn  func(ctx context.Context, id string) error
}

func (m *mockTestUseCase) Create(ctx context.Context, test *domain.Test) (*domain.Test, error) {
	if m.createFn != nil {
		return m.createFn(ctx, test)
	}
	return nil, nil
}

func (m *mockTestUseCase) GetByID(ctx context.Context, id string) (*domain.Test, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockTestUseCase) List(ctx context.Context) ([]domain.Test, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return nil, nil
}

func (m *mockTestUseCase) Update(ctx context.Context, id string, test *domain.Test) (*domain.Test, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, test)
	}
	return nil, nil
}

func (m *mockTestUseCase) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

// --- Helpers ---

func setupTestModuleApp(uc domain.TestUseCase) *fiber.App {
	app := fiber.New()
	handler := rest.NewTestHandler(uc)
	handler.SetupRoutes(app)
	return app
}

type testRequest struct {
	TestID      string `json:"test_id"`
	TestName    string `json:"test_name"`
	Description string `json:"description"`
	Steps       string `json:"steps"`
	Status      string `json:"status"`
	Remarks     string `json:"remarks"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
	IsActive    bool   `json:"is_active"`
}

var sampleTest = domain.Test{
	TestID:      "test-001",
	TestName:    "Login Flow Test",
	Description: "Verify user login works",
	Steps:       "1. Open app\n2. Enter credentials\n3. Submit",
	Status:      "active",
	Remarks:     "Critical path",
	CreatedBy:   "admin",
	CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	UpdatedBy:   "admin",
	UpdatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	IsActive:    true,
}

// --- POST /api/v1/test ---

func TestCreateTest_Success(t *testing.T) {
	mockUC := &mockTestUseCase{
		createFn: func(ctx context.Context, test *domain.Test) (*domain.Test, error) {
			test.TestID = "generated-id"
			return test, nil
		},
	}
	app := setupTestModuleApp(mockUC)

	body, _ := json.Marshal(testRequest{
		TestName:    "New Test",
		Description: "A new test case",
		Steps:       "Step 1\nStep 2",
		Status:      "draft",
		CreatedBy:   "tester",
		IsActive:    true,
	})
	req := httptest.NewRequest("POST", "/api/v1/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte(`"test_name":"New Test"`)) {
		t.Errorf("Expected test_name in response: %s", string(respBody))
	}
}

func TestCreateTest_InvalidPayload(t *testing.T) {
	app := setupTestModuleApp(&mockTestUseCase{})

	req := httptest.NewRequest("POST", "/api/v1/test", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateTest_Error(t *testing.T) {
	mockUC := &mockTestUseCase{
		createFn: func(ctx context.Context, test *domain.Test) (*domain.Test, error) {
			return nil, errors.New("database error")
		},
	}
	app := setupTestModuleApp(mockUC)

	body, _ := json.Marshal(testRequest{TestName: "Fail Test"})
	req := httptest.NewRequest("POST", "/api/v1/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}
}

// --- GET /api/v1/test ---

func TestListTests_Success(t *testing.T) {
	mockUC := &mockTestUseCase{
		listFn: func(ctx context.Context) ([]domain.Test, error) {
			return []domain.Test{sampleTest}, nil
		},
	}
	app := setupTestModuleApp(mockUC)

	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte(`"test_id":"test-001"`)) {
		t.Errorf("Expected test_id in response: %s", string(respBody))
	}
	if !bytes.Contains(respBody, []byte(`"test_name":"Login Flow Test"`)) {
		t.Errorf("Expected test_name in response: %s", string(respBody))
	}
}

func TestListTests_Empty(t *testing.T) {
	mockUC := &mockTestUseCase{
		listFn: func(ctx context.Context) ([]domain.Test, error) {
			return []domain.Test{}, nil
		},
	}
	app := setupTestModuleApp(mockUC)

	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestListTests_Error(t *testing.T) {
	mockUC := &mockTestUseCase{
		listFn: func(ctx context.Context) ([]domain.Test, error) {
			return nil, errors.New("db connection lost")
		},
	}
	app := setupTestModuleApp(mockUC)

	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}
}

// --- GET /api/v1/test/:id ---

func TestGetTestByID_Success(t *testing.T) {
	mockUC := &mockTestUseCase{
		getByIDFn: func(ctx context.Context, id string) (*domain.Test, error) {
			if id == "test-001" {
				return &sampleTest, nil
			}
			return nil, errors.New("not found")
		},
	}
	app := setupTestModuleApp(mockUC)

	req := httptest.NewRequest("GET", "/api/v1/test/test-001", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte(`"test_id":"test-001"`)) {
		t.Errorf("Expected test_id in response: %s", string(respBody))
	}
}

func TestGetTestByID_NotFound(t *testing.T) {
	mockUC := &mockTestUseCase{
		getByIDFn: func(ctx context.Context, id string) (*domain.Test, error) {
			return nil, errors.New("test not found")
		},
	}
	app := setupTestModuleApp(mockUC)

	req := httptest.NewRequest("GET", "/api/v1/test/nonexistent", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte(`"error":"test not found"`)) {
		t.Errorf("Unexpected error response: %s", string(respBody))
	}
}

// --- PUT /api/v1/test/:id ---

func TestUpdateTest_Success(t *testing.T) {
	mockUC := &mockTestUseCase{
		updateFn: func(ctx context.Context, id string, test *domain.Test) (*domain.Test, error) {
			return &domain.Test{
				TestID:      id,
				TestName:    test.TestName,
				Description: test.Description,
				Steps:       test.Steps,
				Status:      test.Status,
				Remarks:     test.Remarks,
				CreatedBy:   "admin",
				CreatedAt:   sampleTest.CreatedAt,
				UpdatedBy:   test.UpdatedBy,
				UpdatedAt:   time.Now(),
				IsActive:    test.IsActive,
			}, nil
		},
	}
	app := setupTestModuleApp(mockUC)

	body, _ := json.Marshal(testRequest{
		TestName:    "Updated Test",
		Description: "Updated description",
		Steps:       "Updated steps",
		Status:      "completed",
		Remarks:     "All passed",
		UpdatedBy:   "tester",
		IsActive:    true,
	})
	req := httptest.NewRequest("PUT", "/api/v1/test/test-001", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte(`"test_name":"Updated Test"`)) {
		t.Errorf("Expected updated test_name: %s", string(respBody))
	}
	if !bytes.Contains(respBody, []byte(`"status":"completed"`)) {
		t.Errorf("Expected updated status: %s", string(respBody))
	}
}

func TestUpdateTest_InvalidPayload(t *testing.T) {
	app := setupTestModuleApp(&mockTestUseCase{})

	req := httptest.NewRequest("PUT", "/api/v1/test/test-001", bytes.NewBufferString(`broken`))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateTest_Error(t *testing.T) {
	mockUC := &mockTestUseCase{
		updateFn: func(ctx context.Context, id string, test *domain.Test) (*domain.Test, error) {
			return nil, errors.New("test not found")
		},
	}
	app := setupTestModuleApp(mockUC)

	body, _ := json.Marshal(testRequest{TestName: "No Such Test"})
	req := httptest.NewRequest("PUT", "/api/v1/test/nonexistent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}
}

// --- DELETE /api/v1/test/:id ---

func TestDeleteTest_Success(t *testing.T) {
	mockUC := &mockTestUseCase{
		deleteFn: func(ctx context.Context, id string) error {
			return nil
		},
	}
	app := setupTestModuleApp(mockUC)

	req := httptest.NewRequest("DELETE", "/api/v1/test/test-001", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte(`"message":"test deleted successfully"`)) {
		t.Errorf("Unexpected response: %s", string(respBody))
	}
}

func TestDeleteTest_Error(t *testing.T) {
	mockUC := &mockTestUseCase{
		deleteFn: func(ctx context.Context, id string) error {
			return errors.New("test not found")
		},
	}
	app := setupTestModuleApp(mockUC)

	req := httptest.NewRequest("DELETE", "/api/v1/test/nonexistent", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}
}
