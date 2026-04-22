package domain

import (
	"context"
	"time"
)

type Test struct {
	TestID      string    `json:"test_id"`
	TestName    string    `json:"test_name"`
	Description string    `json:"description"`
	Steps       string    `json:"steps"`
	Status      string    `json:"status"`
	Remarks     string    `json:"remarks"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedBy   string    `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active"`
}

type TestRepository interface {
	Create(ctx context.Context, test *Test) (*Test, error)
	GetByID(ctx context.Context, id string) (*Test, error)
	List(ctx context.Context) ([]Test, error)
	Update(ctx context.Context, id string, test *Test) (*Test, error)
	Delete(ctx context.Context, id string) error
}

type TestUseCase interface {
	Create(ctx context.Context, test *Test) (*Test, error)
	GetByID(ctx context.Context, id string) (*Test, error)
	List(ctx context.Context) ([]Test, error)
	Update(ctx context.Context, id string, test *Test) (*Test, error)
	Delete(ctx context.Context, id string) error
}
