package http

import (
	"cleo.com/internal/core/port"
	"github.com/gin-gonic/gin"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
}

func NewRouter(
	authService port.AuthService,
	handler HealthMetricParserHandler) (*Router, error) {
	router := gin.Default()
	router.Use(MaxBytesMiddleware(64 * 1024))
	clinicalUser := router.Group("/").Use(RequireClinicalEditor(authService))
	{
		clinicalUser.POST("/parse", handler.Parse)
	}
	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
