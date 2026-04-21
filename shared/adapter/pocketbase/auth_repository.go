package pocketbase

import (
	"context"
	"errors"

	"github.com/pocketbase/pocketbase"
	"boilerplate/services/auth/domain"
)

type authRepository struct {
	pb *pocketbase.PocketBase
}

func NewAuthRepository(pb *pocketbase.PocketBase) domain.AuthRepository {
	return &authRepository{pb: pb}
}

func (r *authRepository) Login(ctx context.Context, username, password string) (*domain.User, string, error) {
	return &domain.User{ID: "stub-id", Username: username}, "stub-jwt-token", nil
}

func (r *authRepository) Signup(ctx context.Context, username, password string) (*domain.User, error) {
	return &domain.User{ID: "new-id", Username: username}, nil
}

func (r *authRepository) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}
	return &domain.User{ID: "valid-id", Username: "validated-user"}, nil
}