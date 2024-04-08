package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/sources/cinemark"
	"github.com/diogovalentte/homarr-iframes/src/sources/linkwarden"
	"github.com/diogovalentte/homarr-iframes/src/sources/vikunja"
)

func IFrameRoutes(group *gin.RouterGroup) {
	group = group.Group("/iframe")
	group.GET("/linkwarden", LinkwardenHandler)
	group.GET("/cinemark", CinemarkHandler)
	group.GET("/vikunja", VikunjaHandler)
}

// @Summary Linkwarden  bookmarks iFrame
// @Description Returns an iFrame with Linkwarden bookmarks.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param collectionId query int false "Get bookmarks only from this collection. You can get the collection ID by going to the collection page. The ID should be on the URL. The ID of the default collection **Unorganized** is 1 because the URL is https://domain.com/collections/1." Example(1)
// @Param theme query string false "Homarr theme, defaults to light." Example(light)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Router /iframe/linkwarden [get]
func LinkwardenHandler(c *gin.Context) {
	l := linkwarden.Linkwarden{}
	err := l.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	l.GetiFrame(c)
}

// @Summary Cinemark Brazil iFrame
// @Description Returns an iFrame with the Cinemark movies in theaters for a city from Brazil.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param city query string true "City to get movies from. First, check if the Cinemark site has a page for this city, if it doesn't, it'll return the page of São Paulo by default. Go to https://cinemark.com.br/rio-de-janeiro/filmes/em-cartaz and select your city. Then grab the city name on the URL." Example(sao-paulo)
// @Param theaters query string false "Thaters' IDs to get movies from. You can find the filter keywords by going to your city page, like https://cinemark.com.br/sao-paulo/filmes/em-cartaz, clicking to filter by theater, and then grabbing the filters in the URL. The filter is the theaters' IDs separated by **%2C**. For example, in the URL https://cinemark.com.br/sao-paulo/filmes/em-cartaz?cinema=716%2C690%2C699 we have the IDs 716, 690, and 699. You have to pass the text `716%2C690%2C699` to the API!" Example(716%2C690%2C699)
// @Param theme query string false "Homarr theme, defaults to light." Example(light)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Router /iframe/cinemark [get]
func CinemarkHandler(c *gin.Context) {
	cin := cinemark.Cinemark{}
	cin.GetiFrame(c)
}

// @Summary Vikunja tasks iFrame
// @Description Returns an iFrame with Vikunja tasks.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param theme query string false "Homarr theme, defaults to light." Example(light)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Router /iframe/vikunja [get]
func VikunjaHandler(c *gin.Context) {
	v := vikunja.Vikunja{}
	err := v.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	v.GetiFrame(c)
}
