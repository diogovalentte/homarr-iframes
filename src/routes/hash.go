package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources/alarms"
	"github.com/diogovalentte/homarr-iframes/src/sources/cinemark"
	"github.com/diogovalentte/homarr-iframes/src/sources/jellyfin"
	"github.com/diogovalentte/homarr-iframes/src/sources/linkwarden"
	"github.com/diogovalentte/homarr-iframes/src/sources/media"
	mediarequets "github.com/diogovalentte/homarr-iframes/src/sources/media-requets"
	uptimekuma "github.com/diogovalentte/homarr-iframes/src/sources/uptime-kuma"
	"github.com/diogovalentte/homarr-iframes/src/sources/vikunja"
)

func HashRoutes(group *gin.RouterGroup) {
	group = group.Group("/hash")
	group.GET("/linkwarden", LinkwardenHashHandler)
	group.GET("/cinemark", CinemarkHashHandler)
	group.GET("/vikunja", VikunjaHashHandler)
	group.GET("/media_releases", MediaReleasesHashHandler)
	group.GET("/media_requests", MediaRequestsHashHandler)
	group.GET("/uptimekuma", UptimeKumaHashHandler)
	group.GET("/alarms", AlarmsHashHandler)
	group.GET("/jellyfin/recently", JellyfinRecentlyHashHandler)
	group.GET("/jellyfin/sessions", JellyfinSessionsHashHandler)
}

// @Summary Get the hash of the Jellyfin sessions
// @Description Get the hash of the Jellyfin active sessions. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param limit query int false "Limits the number of items in the iFrame." Example(20)
// @Param activeWithinSeconds query int false "Only include sessions that have been active within this many seconds. Defaults to 60 if not specified or less than 1." Example(300)
// @Router /hash/jellyfin/sessions [get]
func JellyfinSessionsHashHandler(c *gin.Context) {
	j, err := jellyfin.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	j.GetSessionsHash(c)
}

// @Summary Get the hash of the Jellyfin Recently Added items
// @Description Get the hash of the Jellyfin Recently Added items. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param limit query int false "Limits the number of items in the iFrame." Example(20)
// @Param userId query string false "Jellyfin user ID to get items for. You can get the user ID by going to the admin users page. The ID should be on the URL when clicking on a user. The given user should have access to each library you want data to be fetched from. Defaults to environment JELLYFIN_ADMIN_USER_ID" Example(n6dcfgwiwh1m4c2vhjjm52101vrp01a5)
// @Param parentId query string false "Jellyfin parent/library ID to filter items. You can get the user ID by going to the library page. The ID should be on the URL. The given user in userId should have access to each library you want data to be fetched from." Example(op2xj0l1qejb9g5c5z0f199tigksn1ei)
// @Param includeItemTypes query string false "Filter by media types. Available types can be read here https://api.jellyfin.org/#tag/UserLibrary/operation/GetLatestMedia. Defaults to Movie,Series." Example(Movie,Series)
// @Param queryLimit query int false "Maximum number of items beeing queried from Jellyfin, when fetching Series data each episode counts as one. Defaults to 100." Example(100)
// @Router /hash/jellyfin/recently [get]
func JellyfinRecentlyHashHandler(c *gin.Context) {
	j, err := jellyfin.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	j.GetItemsHash(c)
}

// @Summary Get the hash of the Linkwarden bookmarks
// @Description Get the hash of the Linkwarden bookmarks. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param collectionId query int false "Get bookmarks only from this collection. You can get the collection ID by going to the collection page. The ID should be on the URL. The ID of the default collection **Unorganized** is 1 because the URL is https://domain.com/collections/1." Example(1)
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Router /hash/linkwarden [get]
func LinkwardenHashHandler(c *gin.Context) {
	l, err := linkwarden.New()
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
// @Param theaterIds query string true "The theater IDs to get movies from. It used to be easy to get, but now it's harder. To get it, you need to access the cinemark site, select a theater, open your browser developer console, go to the "Network" tab, filter using the 'onDisplayByTheater' term, and get the theaterId value from the request URL. You have to do it for every theater. Example: 'theaterIds=715, 1222, 4555'" Example(715, 1222, 4555)
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
// @Param project_id query int false "Project ID to get tasks from. You can get it by going to the project page in Vikunja, the project ID should be on the URL. Example project page URL: https://vikunja.com/projects/2, the project ID is 2. Inbox tasks = 1, Favorite tasks = -1." Example(1)
// @Router /hash/vikunja [get]
func VikunjaHashHandler(c *gin.Context) {
	v, err := vikunja.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	v.GetHash(c)
}

// @Summary Get the hash of media releases
// @Description Get the hash of the media releases. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param radarrReleaseType query string false "Filter movies get from Radarr. Can be 'inCinemas', 'physical', or 'digital'. Defaults to 'inCinemas'" Example(physical)
// @Param showUnmonitored query bool false "Specify if show unmonitored media. Defaults to false." Example(true)
// @Router /hash/media_releases [get]
func MediaReleasesHashHandler(c *gin.Context) {
	media.GetHash(c)
}

// @Summary Get the hash of media requests
// @Description Get the hash of the media requests. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param limit query int false "Limits the number of items in the iFrame." Example(5)
// @Param filter query string false "Filters for request status and media status. Available values: all, approved, available, pending, processing, unavailable, failed, deleted, completed, allavaliable (showMedia=true). Defaults to all" Example(all)
// @Param sort query string false "Available values: added, modified, mediaAdded (showMedia=true). Defaults to added" Example(added)
// @Param requestedByOverseerr query string false "If specified, only requests from that particular overseerr user ID will be returned." Example(1)
// @Param requestedByJellyseerr query string false "If specified, only requests from that particular jellyseerr user ID will be returned." Example(1)
// @Param showMedia query string false "If true, shows the requests' media data, not the requests and media data. Defaults to false." Example(true)
// @Router /hash/media_requests [get]
func MediaRequestsHashHandler(c *gin.Context) {
	mediarequets.GetHash(c)
}

// @Summary Get the hash of the Uptime Kuma sites status
// @Description Get the hash of the Uptime Kuma sites status. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param slug query string true "You need to create a status page in Uptime Kuma and select which sites/services this status page will show. While creating the status page, it'll request **you** to create a slug, after creating the status page, provide this slug here. This iFrame will show data only of the sites/services of this specific status page!" Example(uptime-kuma-slug)
// @Router /hash/uptimekuma [get]
func UptimeKumaHashHandler(c *gin.Context) {
	u, err := uptimekuma.New(config.GlobalConfigs.UptimeKumaConfigs.Address, config.GlobalConfigs.UptimeKumaConfigs.InternalAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	u.GetHash(c)
}

// @Summary Get the hash of the alarms
// @Description Get the hash of the alarms. Used by the iFrames to check updates and reload the iframe.
// @Success 200 {object} hashResponse
// @Produce json
// @Param alarms query string true "Alarms to show. Available values: netdata, radarr, lidarr, sonarr, prowlarr, speedtest-tracker, pihole, kavita, kaizoku, changedetectionio, backrest, openarchiver" Example(netdata,radarr,sonarr)
// @Param sort_desc query bool false "Sort alarms in descending order. Defaults to false." Example(false)
// @Param regex_include query bool false "Show only alarms that match or not the regex. Default to true." Example(false)
// @Param changedetectionio_show_viewed query bool false "Show viewed alarms from changedetection.io. Defaults to true." Example(false)
// @Router /hash/alarms [get]
func AlarmsHashHandler(c *gin.Context) {
	a, err := alarms.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	a.GetHash(c)
}
