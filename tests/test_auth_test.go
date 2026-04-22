package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"

	"boilerplate/services/auth/delivery/rest"
	"boilerplate/services/auth/domain"
)

// --- Mock ---

type mockAuthUseCase struct {
	signupFn        func(ctx context.Context, username, password string) (*domain.User, error)
	loginFn         func(ctx context.Context, username, password string) (string, error)
	validateTokenFn func(ctx context.Context, token string) (*domain.User, error)
}

func (m *mockAuthUseCase) Signup(ctx context.Context, username, password string) (*domain.User, error) {
	if m.signupFn != nil {
		return m.signupFn(ctx, username, password)
	}
	return nil, nil
}

func (m *mockAuthUseCase) Login(ctx context.Context, username, password string) (string, error) {
	if m.loginFn != nil {
		return m.loginFn(ctx, username, password)
	}
	return "", nil
}

func (m *mockAuthUseCase) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	if m.validateTokenFn != nil {
		return m.validateTokenFn(ctx, token)
	}
	return nil, nil
}

// --- Helpers ---

func setupAuthApp(uc domain.AuthUseCase) *fiber.App {
	app := fiber.New()
	handler := rest.NewAuthHandler(uc)
	handler.SetupRoutes(app)
	return app
}

type signupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// --- Tests ---

func TestAuthSignup_Success(t *testing.T) {
	mockUC := &mockAuthUseCase{
		signupFn: func(ctx context.Context, username, password string) (*domain.User, error) {
			return &domain.User{ID: "user-1", Username: username}, nil
		},
	}
	app := setupAuthApp(mockUC)

	body, _ := json.Marshal(signupRequest{Username: "newuser", Password: "password123"})
	req := httptest.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte(`"message":"Signup successful"`)) {
		t.Errorf("Unexpected response: %s", string(respBody))
	}
	if !bytes.Contains(respBody, []byte(`"Username":"newuser"`)) {
		t.Errorf("Expected username in response: %s", string(respBody))
	}
}

func TestAuthSignup_InvalidPayload(t *testing.T) {
	app := setupAuthApp(&mockAuthUseCase{})

	req := httptest.NewRequest("POST", "/api/v1/signup", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestAuthSignup_Error(t *testing.T) {
	mockUC := &mockAuthUseCase{
		signupFn: func(ctx context.Context, username, password string) (*domain.User, error) {
			return nil, errors.New("username already exists")
		},
	}
	app := setupAuthApp(mockUC)

	body, _ := json.Marshal(signupRequest{Username: "existing", Password: "password123"})
	req := httptest.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte(`"error":"username already exists"`)) {
		t.Errorf("Unexpected error response: %s", string(respBody))
	}
}

func TestAuthLogin_Success(t *testing.T) {
	mockUC := &mockAuthUseCase{
		loginFn: func(ctx context.Context, username, password string) (string, error) {
			return "jwt-token-abc", nil
		},
	}
	app := setupAuthApp(mockUC)

	body, _ := json.Marshal(loginRequest{Username: "testuser", Password: "testpass"})
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte(`"token":"jwt-token-abc"`)) {
		t.Errorf("Expected token in response: %s", string(respBody))
	}
	if !bytes.Contains(respBody, []byte(`"message":"Login successful"`)) {
		t.Errorf("Expected success message: %s", string(respBody))
	}
}

func TestAuthLogin_InvalidPayload(t *testing.T) {
	app := setupAuthApp(&mockAuthUseCase{})

	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(`not-json`))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestAuthLogin_InvalidCredentials(t *testing.T) {
	mockUC := &mockAuthUseCase{
		loginFn: func(ctx context.Context, username, password string) (string, error) {
			return "", errors.New("invalid credentials")
		},
	}
	app := setupAuthApp(mockUC)

	body, _ := json.Marshal(loginRequest{Username: "wrong", Password: "wrong"})
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte(`"error":"Invalid credentials"`)) {
		t.Errorf("Unexpected error response: %s", string(respBody))
	}
}

func TestAuthLogout_Success(t *testing.T) {
	app := setupAuthApp(&mockAuthUseCase{})

	req := httptest.NewRequest("POST", "/api/v1/logout", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte(`"message":"Logged out successfully"`)) {
		t.Errorf("Unexpected response: %s", string(respBody))
	}
}
