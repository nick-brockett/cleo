package http

import (
	"net/http"

	"cleo.com/internal/core/port"
	"github.com/gin-gonic/gin"
)

func MaxBytesMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

func RequireClinicalEditor(s port.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		hasRole, err := s.HasRole(c, "CLINICAL-EDITOR")
		if err != nil || !hasRole {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
