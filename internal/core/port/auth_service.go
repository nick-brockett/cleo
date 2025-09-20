package port

import "github.com/gin-gonic/gin"

type AuthService interface {
	HasRole(c *gin.Context, role string) (bool, error)
}
