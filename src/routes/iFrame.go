package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources/alarms"
	"github.com/diogovalentte/homarr-iframes/src/sources/cinemark"
	"github.com/diogovalentte/homarr-iframes/src/sources/linkwarden"
	"github.com/diogovalentte/homarr-iframes/src/sources/media"
	mediarequets "github.com/diogovalentte/homarr-iframes/src/sources/media-requets"
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
	group.GET("/media_requests", MediaRequestsiFrameHandler)
	group.GET("/uptimekuma", UptimeKumaiFrameHandler)
	group.GET("/alarms", AlarmsiFrameHandler)
	group.GET("/netdata", NetdataiFrameHandler)
}

// @Summary Linkwarden  bookmarks iFrame
// @Description Returns an iFrame with Linkwarden bookmarks.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param collectionId query int false "Get bookmarks only from this collection. You can get the collection ID by going to the collection page. The ID should be on the URL. The ID of the default collection **Unorganized** is 1 because the URL is https://domain.com/collections/1." Example(1)
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload." Example(https://sub.domain.com)
// @Param background_position query string false "Background position of each bookmark card. Use '%25' in place of '%', like '50%25 47.2%25' to get '50% 47.2%'. Defaults to 50% 47.2%." Example(top)
// @Param background_size query string false "Background size of each bookmark card. Use '%25' in place of '%'. Defaults to cover." Example(cover)
// @Param background_filter query string false "Background filter of each bookmark card. Use '%25' in place of '%'. Defaults to brightness(0.3)." Example(blur(5px))
// @Router /iframe/linkwarden [get]
func LinkwardeniFrameHandler(c *gin.Context) {
	l, err := linkwarden.New()
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
// @Param exclude_project_ids query string false "Project IDs to NOT get tasks from. You can get it by going to the project page in Vikunja, the project ID should be on the URL. Example project page URL: https://vikunja.com/projects/2, the project ID is 2. Inbox tasks = 1, Favorite tasks = -1." Example(1,5,7)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload. Also used by the button to set the task done, if not provided, the button will not appear (the button doesn't appear in repeating tasks.)" Example(https://sub.domain.com)
// @Param showCreated query bool false "Shows the tasks' created date. Defaults to true." Example(false)
// @Param showDue query bool false "Shows the tasks' due/end date and repeating dates. Defaults to true." Example(false)
// @Param showPriority query bool false "Shows the tasks' priority. Defaults to true." Example(false)
// @Param showProject query bool false "Shows the tasks' project. Defaults to true." Example(false)
// @Param showFavoriteIcon query bool false "Shows a start icon in favorite tasks. Defaults to true." Example(false)
// @Param showLabels query bool false "Shows the tasks' labels. Defaults to true." Example(false)
// @Param background_position query string false "Background position of each task card. Use '%25' in place of '%', like '50%25 47.2%25' to get '50% 47.2%'. Defaults to 50% 49.5%." Example(top)
// @Param background_size query string false "Background size of each task card. Use '%25' in place of '%'. Defaults to 105%." Example(105%25)
// @Param background_filter query string false "Background filter of each task card. Use '%25' in place of '%'. Defaults to brightness(0.3)." Example(blur(5px))
// @Router /iframe/vikunja [get]
func VikunjaiFrameHandler(c *gin.Context) {
	v, err := vikunja.New()
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
	taskIDStr := c.Query("taskId")
	if taskIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "taskId is required"})
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "taskId must be an integer"})
	}

	v, err := vikunja.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	err = v.SetTaskDone(taskID)
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
	c.String(http.StatusMovedPermanently, "Overseerr iFrame was removed. It's now implemented in the media requests iFrame. Please consult the media requests iFrame documentation.")
}

// @Summary Media Releases
// @Description Returns an iFrame with the media releases of today. The media releases are from Radarr/Sonarr/Lidarr.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload. Also used by the button to set the task done, if not provided, the button will not appear." Example(https://sub.domain.com)
// @Param radarrReleaseType query string false "Filter movies get from Radarr. Can be 'inCinemas', 'physical', 'digital', or multiple separated by comma. Defaults to 'inCinemas,physical,digital'" Example(inCinemas,digital)
// @Param showUnmonitored query bool false "Specify if show unmonitored media. Defaults to false." Example(true)
// @Param showEpisodesHour query bool false "Specify if show the episodes' (Sonarr) release hour and minute. Defaults to true." Example(false)
// @Router /iframe/media_releases [get]
func MediaReleasesiFrameHandler(c *gin.Context) {
	media.GetiFrame(c)
}

// @Summary Overseerr and Jellyseerr Media Requests
// @Description Returns an iFrame with Overseerr and Jellyseerr media requests list. Returns all requests if the user's API token has the ADMIN or MANAGE_REQUESTS permissions. Otherwise, only the logged-in user's requests are returned.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload. Also used by the button to set the task done, if not provided, the button will not appear." Example(https://sub.domain.com)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Param filter query string false "Filters for request status and media status. Available values: all, approved, available, pending, processing, unavailable, failed, deleted, completed. Defaults to all" Example(all)
// @Param sort query string false "Available values: added, modified. Defaults to added" Example(added)
// @Param requestedByOverseerr query string false "If specified, only requests from that particular overseerr user ID will be returned." Example(1)
// @Param requestedByJellyseerr query string false "If specified, only requests from that particular jellyseerr user ID will be returned." Example(1)
// @Router /iframe/media_requests [get]
func MediaRequestsiFrameHandler(c *gin.Context) {
	mediarequets.GetiFrame(c)
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
	u, err := uptimekuma.New(config.GlobalConfigs.UptimeKumaConfigs.Address, config.GlobalConfigs.UptimeKumaConfigs.InternalAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	u.GetiFrame(c)
}

// @Summary Alarms iFrame
// @Description Returns an iFrame with alarms from multiple sources.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Param theme query string false "Homarr theme, defaults to light. If it's different from your Homarr theme, the background turns white" Example(light)
// @Param api_url query string true "API URL used by your browser. Use by the iFrames to check any update, if there is an update, the iFrame reloads. If not specified, the iFrames will never try to reload." Example(https://sub.domain.com)
// @Param alarms query string true "Alarms to show. Available values: netdata, radarr, lidarr, sonarr, prowlarr, speedtest-tracker, pihole, kavita, kaizoku, changedetectionio, backrest" Example(netdata,radarr,sonarr)
// @Param sort_desc query bool false "Sort alarms in descending order. Defaults to false." Example(false)
// @Param regex_include query bool false "Show only alarms that match or not the regex. Default to true." Example(false)
// @Param changedetectionio_show_viewed query bool false "Show viewed alarms from changedetection.io. Defaults to true." Example(false)
// @Router /iframe/alarms [get]
func AlarmsiFrameHandler(c *gin.Context) {
	a, err := alarms.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	a.GetiFrame(c)
}

// @Summary Netdata iFrame
// @Description Returns a message saying that this iFrame is not implemented anymore.
// @Success 200 {string} string "HTML content"
// @Produce html
// @Router /iframe/netdata [get]
func NetdataiFrameHandler(c *gin.Context) {
	c.String(http.StatusMovedPermanently, "Netdata iFrame was removed. It's now implemented in the alarms iFrame. Please consult the alarms iFrame documentation.")
}

type messsageResponse struct {
	Message string `json:"message"`
}
