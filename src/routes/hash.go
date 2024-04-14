package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/sources/cinemark"
	"github.com/diogovalentte/homarr-iframes/src/sources/linkwarden"
	"github.com/diogovalentte/homarr-iframes/src/sources/vikunja"
)

func HashRoutes(group *gin.RouterGroup) {
	group = group.Group("/hash")
	group.GET("/linkwarden", LinkwardenHashHandler)
	group.GET("/cinemark", CinemarkHashHandler)
	group.GET("/vikunja", VikunjaHashHandler)
}

// @Summary Get the hash of the Linkwarden bookmarks
// @Description Get the hash of the Linkwarden bookmarks. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param collectionId query int false "Get bookmarks only from this collection. You can get the collection ID by going to the collection page. The ID should be on the URL. The ID of the default collection **Unorganized** is 1 because the URL is https://domain.com/collections/1." Example(1)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Router /hash/linkwarden [get]
func LinkwardenHashHandler(c *gin.Context) {
	l := linkwarden.Linkwarden{}
	err := l.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	l.GetHash(c)
}

// @Summary Get the hash of the Cinemark movies
// @Description Get the hash of the Cinemark movies. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param city query string true "City to get movies from. First, check if the Cinemark site has a page for this city, if it doesn't, it'll return the page of São Paulo by default. Go to https://cinemark.com.br/rio-de-janeiro/filmes/em-cartaz and select your city. Then grab the city name on the URL." Example(sao-paulo)
// @Param theaters query string false "Thaters' IDs to get movies from. You can find the filter keywords by going to your city page, like https://cinemark.com.br/sao-paulo/filmes/em-cartaz, clicking to filter by theater, and then grabbing the filters in the URL. The filter is the theaters' IDs separated by **%2C**. For example, in the URL https://cinemark.com.br/sao-paulo/filmes/em-cartaz?cinema=716%2C690%2C699 we have the IDs 716, 690, and 699. You have to pass the text `716%2C690%2C699` to the API!" Example(716%2C690%2C699)
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
	v := vikunja.Vikunja{}
	err := v.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	v.GetHash(c)
}
