package jwt

import (
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func GenerateToken(userID string, username string, secret string) (string, error) {
	if secret == "" {
		return "", errors.New("jwt secret is empty")
	}

	token, err := jwt.NewBuilder().
		Subject(userID).
		Claim("username", username).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(24 * time.Hour)).
		Build()

	if err != nil {
		return "", err
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256(), []byte(secret)))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func ValidateToken(tokenString string, secret string) (string, string, error) {
	if secret == "" {
		return "", "", errors.New("jwt secret is empty")
	}

	token, err := jwt.Parse([]byte(tokenString), jwt.WithKey(jwa.HS256(), []byte(secret)), jwt.WithValidate(true))
	if err != nil {
		return "", "", err
	}

	userID, ok := token.Subject()
	if !ok {
		return "", "", errors.New("missing subject claim")
	}

	var username string
	if err := token.Get("username", &username); err != nil {
		return "", "", errors.New("missing or invalid username claim")
	}

	return userID, username, nil
}