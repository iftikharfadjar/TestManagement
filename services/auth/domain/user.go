package domain

import "context"

type User struct {
	ID       string
	Username string
}

type AuthRepository interface {
	Login(ctx context.Context, username, password string) (*User, string, error)
	Signup(ctx context.Context, username, password string) (*User, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
}

type AuthUseCase interface {
	Login(ctx context.Context, username, password string) (string, error)
	Signup(ctx context.Context, username, password string) (*User, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
}