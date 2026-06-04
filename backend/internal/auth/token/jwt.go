package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret []byte
}

func New(secret string) *Service {
	return &Service{
		secret: []byte(secret),
	}
}

func (s *Service) GenerateAccessToken(
	userID string,
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
