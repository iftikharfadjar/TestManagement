package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"boilerplate/services/auth/domain"
)

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

func setupTestApp(uc domain.AuthUseCase) *fiber.App {
	app := fiber.New()
	handler := NewAuthHandler(uc)
	handler.SetupRoutes(app)
	return app
}

func TestSignup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUC := &mockAuthUseCase{
			signupFn: func(ctx context.Context, username, password string) (*domain.User, error) {
				return &domain.User{ID: "123", Username: username}, nil
			},
		}

		app := setupTestApp(mockUC)

		body, _ := json.Marshal(SignupRequest{Username: "testuser", Password: "testpassword"})
		req := httptest.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		respBody, _ := io.ReadAll(resp.Body)
		if !bytes.Contains(respBody, []byte(`"user":{"ID":"123","Username":"testuser"}`)) {
			t.Errorf("Unexpected response body: %s", string(respBody))
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockUC := &mockAuthUseCase{
			signupFn: func(ctx context.Context, username, password string) (*domain.User, error) {
				return nil, errors.New("username already exists")
			},
		}

		app := setupTestApp(mockUC)

		body, _ := json.Marshal(SignupRequest{Username: "testuser", Password: "testpassword"})
		req := httptest.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}
	})
}

func TestLogin(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUC := &mockAuthUseCase{
			loginFn: func(ctx context.Context, username, password string) (string, error) {
				return "valid-test-token", nil
			},
		}

		app := setupTestApp(mockUC)

		body, _ := json.Marshal(LoginRequest{Username: "testuser", Password: "testpassword"})
		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		respBody, _ := io.ReadAll(resp.Body)
		if !bytes.Contains(respBody, []byte(`"token":"valid-test-token"`)) {
			t.Errorf("Unexpected response body: %s", string(respBody))
		}
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		mockUC := &mockAuthUseCase{
			loginFn: func(ctx context.Context, username, password string) (string, error) {
				return "", errors.New("invalid credentials")
			},
		}

		app := setupTestApp(mockUC)

		body, _ := json.Marshal(LoginRequest{Username: "wrong", Password: "wrong"})
		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}

func TestLogout(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		app := setupTestApp(&mockAuthUseCase{})

		req := httptest.NewRequest("POST", "/api/v1/logout", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		
		respBody, _ := io.ReadAll(resp.Body)
		if !bytes.Contains(respBody, []byte(`"message":"Logged out successfully"`)) {
			t.Errorf("Unexpected response body: %s", string(respBody))
		}
	})
}