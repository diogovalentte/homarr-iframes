package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources/cinemark"
	"github.com/diogovalentte/homarr-iframes/src/sources/linkwarden"
	uptimekuma "github.com/diogovalentte/homarr-iframes/src/sources/uptime-kuma"
	"github.com/diogovalentte/homarr-iframes/src/sources/vikunja"
)

func HashRoutes(group *gin.RouterGroup) {
	group = group.Group("/hash")
	group.GET("/linkwarden", LinkwardenHashHandler)
	group.GET("/cinemark", CinemarkHashHandler)
	group.GET("/vikunja", VikunjaHashHandler)
	group.GET("/uptimekuma", UptimeKumaHashHandler)
}

// @Summary Get the hash of the Linkwarden bookmarks
// @Description Get the hash of the Linkwarden bookmarks. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param collectionId query int false "Get bookmarks only from this collection. You can get the collection ID by going to the collection page. The ID should be on the URL. The ID of the default collection **Unorganized** is 1 because the URL is https://domain.com/collections/1." Example(1)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Router /hash/linkwarden [get]
func LinkwardenHashHandler(c *gin.Context) {
	l, err := linkwarden.New(config.GlobalConfigs.Linkwarden.Address, config.GlobalConfigs.Linkwarden.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	l.GetHash(c)
}

// @Summary Get the hash of the Cinemark movies
// @Description Get the hash of the Cinemark movies. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param theaterIs query string true "The theater IDs to get movies from. It used to be easy to get, but now it's harder. To get it, you need to access the cinemark site, select a theater, open your browser developer console, go to the "Network" tab, filter using the 'onDisplayByTheater' term, and get the theaterId value from the request URL. You have to do it for every theater. Example: 'theaterIds=715, 1222, 4555'" Example(715, 1222, 4555)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Router /hash/cinemark [get]
func CinemarkHashHandler(c *gin.Context) {
	cin := cinemark.Cinemark{}
	cin.GetHash(c)
}

type hashResponse struct {
	Hash string `json:"hash"`
}

// @Summary Get the hash of the Vikunja tasks
// @Description Get the hash of the Vikunja tasks. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Router /hash/vikunja [get]
func VikunjaHashHandler(c *gin.Context) {
	v, err := vikunja.New(config.GlobalConfigs.Vikunja.Address, config.GlobalConfigs.Vikunja.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	v.GetHash(c)
}

// @Summary Get the hash of the Uptime Kuma sites status
// @Description Get the hash of the Uptime Kuma sites status. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param slug query string true "You need to create a status page in Uptime Kuma and select which sites/services this status page will show. While creating the status page, it'll request **you** to create a slug, after creating the status page, provide this slug here. This iFrame will show data only of the sites/services of this specific status page!" Example(uptime-kuma-slug)
// @Router /hash/uptimekuma [get]
func UptimeKumaHashHandler(c *gin.Context) {
	u, err := uptimekuma.New(config.GlobalConfigs.UptimeKumaConfigs.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	u.GetHash(c)
}
