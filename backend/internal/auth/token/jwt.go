package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	secret []byte
}

type GenerationTokens struct {
	Access  string
	Refresh string
}

func New(secret string) *Service {
	return &Service{
		secret: []byte(secret),
	}
}

func GenerateSessionID() uuid.UUID {
	return uuid.New()
}

func GenerateRefreshSecret() (string, error) {

	b := make([]byte, 64)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *Service) GenerateAccessToken(
	userID uuid.UUID,
) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userID,
			"exp": time.Now().
				Add(1 * time.Hour).
				Unix(),
		},
	)

	return token.SignedString(s.secret)
}

func (s *Service) ParseAccessToken(
	tokenString string,
) (string, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return s.secret, nil
		},
	)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token")
	}

	return userID, nil
}

func ParseRefreshToken(
	token string,
) (uuid.UUID, string, error) {

	parts := strings.Split(token, ".")

	if len(parts) != 2 {
		return uuid.Nil, "", errors.New(
			"invalid refresh token",
		)
	}

	sessionID, err := uuid.Parse(
		parts[0],
	)

	if err != nil {
		return uuid.Nil, "", err
	}

	return sessionID, parts[1], nil
}
