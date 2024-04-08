package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/sources/linkwarden"
)

// LinksRoutes registers the links routes
func IFrameRoutes(group *gin.RouterGroup) {
	group = group.Group("/iframe")
	group.GET("/linkwarden", LinkwardenHandler)
}

func LinkwardenHandler(c *gin.Context) {
	l := linkwarden.Linkwarden{}
	err := l.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"response": err.Error()})
	}
	l.GetiFrame(c)
}
