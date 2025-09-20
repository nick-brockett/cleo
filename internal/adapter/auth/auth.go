package auth

import (
	"errors"
	"fmt"
	"strings"

	"cleo.com/internal/core/port"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type Config struct {
	APISecret string `env:"API_SECRET, default=a-string-secret-at-least-256-bits-long"`
}

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Service struct {
	logger    *logrus.Logger
	APISecret string
}

func NewService(logger *logrus.Logger, config Config) port.AuthService {
	return &Service{
		logger:    logger,
		APISecret: config.APISecret,
	}
}

func (s *Service) HasRole(c *gin.Context, role string) (bool, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return false, errors.New("missing authorization header")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		s.logger.Error("invalid authorization header")
		return false, errors.New("invalid authorization header")
	}
	tokenStr := parts[1]

	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.APISecret), nil
	})
	if err != nil || !token.Valid {
		return false, errors.New("invalid token")
	}

	return claims.Role == role, nil
}
