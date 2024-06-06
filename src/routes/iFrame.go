package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources/cinemark"
	"github.com/diogovalentte/homarr-iframes/src/sources/linkwarden"
	"github.com/diogovalentte/homarr-iframes/src/sources/media"
	"github.com/diogovalentte/homarr-iframes/src/sources/overseerr"
	uptimekuma "github.com/diogovalentte/homarr-iframes/src/sources/uptime-kuma"
	"github.com/diogovalentte/homarr-iframes/src/sources/vikunja"
)

func IFrameRoutes(group *gin.RouterGroup) {
	group = group.Group("/iframe")
	group.GET("/linkwarden", LinkwardeniFrameHandler)
	group.GET("/cinemark", CinemarkiFrameHandler)
	group.GET("/vikunja", VikunjaiFrameHandler)
	group.PATCH("/vikunja/set_task_done", VikunjaSetTaskDoneHandler)
	group.GET("/overseerr", OverseerriFrameHandler)
	group.GET("/media_releases", MediaReleasesiFrameHandler)
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
	l, err := linkwarden.New(config.GlobalConfigs.Linkwarden.Address, config.GlobalConfigs.Linkwarden.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	l.GetiFrame(c)
}

// @Summary Cinemark Brazil iFrame
// @Description Returns an iFrame with the on display movies in specific Cinemark theaters. I recommend you to get the movies from the theaters of your city.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param theaterIds query string true "The theater IDs to get movies from. It used to be easy to get, but now it's harder. To get it, you need to access the cinemark site, select a theater, open your browser developer console, go to the "Network" tab, filter using the 'onDisplayByTheater' term, and get the theaterId value from the request URL. You have to do it for every theater. Example: 'theaterIds=715, 1222, 4555'" Example(715, 1222, 4555)
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload." Example(https://sub.domain.com)
// @Router /iframe/cinemark [get]
func CinemarkiFrameHandler(c *gin.Context) {
	cin := cinemark.Cinemark{}
	cin.GetiFrame(c)
}

// @Summary Vikunja tasks iFrame
// @Description Returns an iFrame with not done Vikunja tasks. Uses a custom sort/order: due date (asc); end date (asc); priority (desc); created date (desc). When the due/end date is today, the date color will be orange, if it's past due, the date color will be red.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Param project_id query int false "Project ID to get tasks from. You can get it by going to the project page in Vikunja, the project ID should be on the URL. Example project page URL: https://vikunja.com/projects/2, the project ID is 2. Inbox tasks = 1, Favorite tasks = -1." Example(1)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload. Also used by the button to set the task done, if not provided, the button will not appear (the button doesn't appear in repeating tasks.)" Example(https://sub.domain.com)
// @Param showCreated query bool false "Shows the tasks' created date. Defaults to true." Example(false)
// @Param showDue query bool false "Shows the tasks' due/end date and repeating dates. Defaults to true." Example(false)
// @Param showPriority query bool false "Shows the tasks' priority. Defaults to true." Example(false)
// @Param showProject query bool false "Shows the tasks' project. Defaults to true." Example(false)
// @Param showFavoriteIcon query bool false "Shows a start icon in favorite tasks. Defaults to true." Example(false)
// @Router /iframe/vikunja [get]
func VikunjaiFrameHandler(c *gin.Context) {
	v, err := vikunja.New(config.GlobalConfigs.Vikunja.Address, config.GlobalConfigs.Vikunja.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	v.GetiFrame(c)
}

// @Summary Set Vikunja task done
// @Description Set a Vikunja task as done.
// @Success 200 {object} messsageResponse "Task done"
// @Produce json
// @Param taskId query int true "The task ID." Example(1)
// @Router /iframe/vikunja/set_task_done [patch]
func VikunjaSetTaskDoneHandler(c *gin.Context) {
	taskIdStr := c.Query("taskId")
	if taskIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "taskId is required"})
		return
	}

	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "taskId must be an integer"})
	}

	v, err := vikunja.New(config.GlobalConfigs.Vikunja.Address, config.GlobalConfigs.Vikunja.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	err = v.SetTaskDone(taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task done"})
}

// @Summary Overseerr Media Requests
// @Description Returns an iFrame with Overseerr media requests list. Returns all requests if the user's API token has the ADMIN or MANAGE_REQUESTS permissions. Otherwise, only the logged-in user's requests are returned.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload. Also used by the button to set the task done, if not provided, the button will not appear." Example(https://sub.domain.com)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Param filter query string false "Available values : all, approved, available, pending, processing, unavailable, failed" Example(all)
// @Param sort query string false "Available values : added, modified. Defaults to added" Example(added)
// @Param requestedBy query string false "If specified, only requests from that particular user ID will be returned." Example(1)
// @Router /iframe/overseerr [get]
func OverseerriFrameHandler(c *gin.Context) {
	o, err := overseerr.New(config.GlobalConfigs.Overseerr.Address, config.GlobalConfigs.Overseerr.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	o.GetiFrame(c)
}

// @Summary Media Releases
// @Description Returns an iFrame with the media releases of today. The media releases are from Radarr/Sonarr.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload. Also used by the button to set the task done, if not provided, the button will not appear." Example(https://sub.domain.com)
// @Param radarrReleaseType query string false "Filter movies get from Radarr. Can be 'inCinemas', 'physical', or 'digital'. Defaults to 'inCinemas'" Example(physical)
// @Param showUnmonitored query bool false "Specify if show unmonitored media. Defaults to false." Example(true)
// @Router /iframe/media_releases [get]
func MediaReleasesiFrameHandler(c *gin.Context) {
	media.GetiFrame(c)
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
	u, err := uptimekuma.New(config.GlobalConfigs.UptimeKumaConfigs.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	u.GetiFrame(c)
}

type messsageResponse struct {
	Message string `json:"message"`
}
