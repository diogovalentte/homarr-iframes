package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/sources/cinemark"
	"github.com/diogovalentte/homarr-iframes/src/sources/linkwarden"
	"github.com/diogovalentte/homarr-iframes/src/sources/vikunja"
)

// LinksRoutes registers the links routes
func IFrameRoutes(group *gin.RouterGroup) {
	group = group.Group("/iframe")
	group.GET("/linkwarden", LinkwardenHandler)
	group.GET("/cinemark", CinemarkHandler)
	group.GET("/vikunja", VikunjaHandler)
}

func LinkwardenHandler(c *gin.Context) {
	l := linkwarden.Linkwarden{}
	err := l.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	l.GetiFrame(c)
}

func CinemarkHandler(c *gin.Context) {
	cin := cinemark.Cinemark{}
	cin.GetiFrame(c)
}

func VikunjaHandler(c *gin.Context) {
	v := vikunja.Vikunja{}
	err := v.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	v.GetiFrame(c)
}
