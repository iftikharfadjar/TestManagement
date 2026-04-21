package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"boilerplate/services/auth/domain"
	"boilerplate/shared/adapter/sqlite_adapter/sqlc"
	"boilerplate/shared/jwt"
)

type authRepository struct {
	db        *sql.DB
	queries   *sqlc.Queries
	jwtSecret string
}

func NewAuthRepository(db *sql.DB, jwtSecret string) domain.AuthRepository {
	return &authRepository{
		db:        db,
		queries:   sqlc.New(db),
		jwtSecret: jwtSecret,
	}
}

func (r *authRepository) Login(ctx context.Context, username, password string) (*domain.User, string, error) {
	sqlcUser, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", errors.New("invalid credentials")
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(sqlcUser.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(sqlcUser.ID, sqlcUser.Username, r.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return &domain.User{
		ID:       sqlcUser.ID,
		Username: sqlcUser.Username,
	}, token, nil
}

func (r *authRepository) Signup(ctx context.Context, username, password string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()

	err = r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           id,
		Username:     username,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:       id,
		Username: username,
	}, nil
}

func (r *authRepository) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	userID, username, err := jwt.ValidateToken(token, r.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:       userID,
		Username: username,
	}, nil
}