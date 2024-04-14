// Package api implements the API routes and groups
package api

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/diogovalentte/homarr-iframes/docs"
	"github.com/diogovalentte/homarr-iframes/src/routes"
)

// SetupRouter sets up the API routes
func SetupRouter() *gin.Engine {
	router := gin.Default()

	docs.SwaggerInfo.Title = "Homarr iFrames API"
	docs.SwaggerInfo.Description = "iFrames of many applications to use in Homarr"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/v1"

	v1 := router.Group("/v1")
	{
		routes.HealthCheckRoute(v1)
	}
	{
		routes.IFrameRoutes(v1)
	}
	{
		routes.HashRoutes(v1)
	}

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
