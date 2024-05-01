package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/sources/cinemark"
	"github.com/diogovalentte/homarr-iframes/src/sources/linkwarden"
	uptimekuma "github.com/diogovalentte/homarr-iframes/src/sources/uptime-kuma"
	"github.com/diogovalentte/homarr-iframes/src/sources/vikunja"
)

func IFrameRoutes(group *gin.RouterGroup) {
	group = group.Group("/iframe")
	group.GET("/linkwarden", LinkwardeniFrameHandler)
	group.GET("/cinemark", CinemarkiFrameHandler)
	group.GET("/vikunja", VikunjaiFrameHandler)
	group.PATCH("/vikunja/set_task_done", VikunjaSetTaskDoneHandler)
	group.GET("/uptimekuma", UptimeKumaiFrameHandler)
}

// @Summary Linkwarden  bookmarks iFrame
// @Description Returns an iFrame with Linkwarden bookmarks.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param collectionId query int false "Get bookmarks only from this collection. You can get the collection ID by going to the collection page. The ID should be on the URL. The ID of the default collection **Unorganized** is 1 because the URL is https://domain.com/collections/1." Example(1)
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload." Example(https://sub.domain.com)
// @Router /iframe/linkwarden [get]
func LinkwardeniFrameHandler(c *gin.Context) {
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
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload." Example(https://sub.domain.com)
// @Router /iframe/cinemark [get]
func CinemarkiFrameHandler(c *gin.Context) {
	cin := cinemark.Cinemark{}
	cin.GetiFrame(c)
}

// @Summary Vikunja tasks iFrame
// @Description Returns an iFrame with Vikunja tasks.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload. Also used by the button to set the task done, if not provided, the button will not appear." Example(https://sub.domain.com)
// @Param showCreated query bool false "Shows the tasks' created date. Defaults to true." Example(false)
// @Param showDue query bool false "Shows the tasks' due/end date and repeating dates. Defaults to true." Example(false)
// @Param showPriority query bool false "Shows the tasks' priority. Defaults to true." Example(false)
// @Param showProject query bool false "Shows the tasks' project. Defaults to true." Example(false)
// @Router /iframe/vikunja [get]
func VikunjaiFrameHandler(c *gin.Context) {
	v := vikunja.Vikunja{}
	err := v.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	v.GetiFrame(c)
}

// @Summary Set Vikunja task done
// @Description Set a Vikunja task as done.
// @Success 200 {object} messsageResponse "Task done"
// @Produce json
// @Param task_id query int true "The task ID." Example(1)
// @Router /iframe/vikunja/set_task_done [patch]
func VikunjaSetTaskDoneHandler(c *gin.Context) {
	taskIdStr := c.Query("taskId")
	if taskIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "task_id is required"})
		return
	}

	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "task_id must be an integer"})
	}

	v := vikunja.Vikunja{}
	err = v.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	err = v.SetTaskDone(taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task done"})
}

// @Summary Uptime Kuma iFrame
// @Description Returns an iFrame with Uptime Kuma sites overview.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param slug query string true "You need to create a status page in Uptime Kuma and select which sites/services this status page will show. While creating the status page, it'll request **you** to create a slug, after creating the status page, provide this slug here. This iFrame will show data only of the sites/services of this specific status page!" Example(uptime-kuma-slug)
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload." Example(https://sub.domain.com)
// @Param showTitle query bool false "Show the title 'Uptime Kuma' on the iFrame." Example(true)
// @Param orientation query string false "Orientation of the containers, defaults to horizontal." Example(vertical)
// @Router /iframe/uptimekuma [get]
func UptimeKumaiFrameHandler(c *gin.Context) {
	u := uptimekuma.UptimeKuma{}
	err := u.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	u.GetiFrame(c)
}

type messsageResponse struct {
	Message string `json:"message"`
}
