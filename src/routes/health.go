// Package routes implements the API routes
package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Health check route
// @Description Returns status OK
// @Success 200 {string} string OK
// @Produce plain
// @Router /health [get]
func HealthCheckRoute(group *gin.RouterGroup) {
	group.GET("/health", healthCheck)
}

func healthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
