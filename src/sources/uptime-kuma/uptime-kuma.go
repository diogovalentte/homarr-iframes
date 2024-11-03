package uptimekuma

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/sources"
)

var u *UptimeKuma

// UptimeKuma is the UptimeKuma source
type UptimeKuma struct {
	Address         string
	InternalAddress string
}

func New(address, internalAddress string) (*UptimeKuma, error) {
	if u != nil {
		return u, nil
	}

	newU := &UptimeKuma{}
	err := newU.Init(address, internalAddress)
	if err != nil {
		return nil, err
	}

	u = newU

	return u, nil
}

// Init sets the UptimeKuma properties from the configs
func (u *UptimeKuma) Init(address, internalAddress string) error {
	if address == "" {
		return fmt.Errorf("UPTIMEKUMA_ADDRESS variable should be set")
	}

	u.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		u.InternalAddress = u.Address
	} else {
		u.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}

	return nil
}

// GetiFrame returns the iFrame for the UptimeKuma source
func (u *UptimeKuma) GetiFrame(c *gin.Context) {
	theme := c.Query("theme")
	if theme == "" {
		theme = "light"
	} else if theme != "dark" && theme != "light" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "theme must be 'dark' or 'light'"})
		return
	}

	slug := c.Query("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "slug must be provided"})
		return
	}

	var err error
	apiURL := c.Query("api_url")
	if apiURL != "" {
		_, err = url.ParseRequestURI(apiURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "api_url must be a valid URL like 'http://192.168.1.46:8080' or 'https://sub.domain.com'"})
			return
		}
	}

	showTitleStr := c.Query("showTitle")
	showTitle := true
	if showTitleStr != "" {
		showTitle, err = strconv.ParseBool(showTitleStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "title must be a boolean (true or false)"})
			return
		}
	}

	containersDisplay := c.Query("orientation")
	if containersDisplay == "" {
		containersDisplay = "horizontal"
	} else if containersDisplay != "horizontal" && containersDisplay != "vertical" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "orientation must be 'horizontal' or 'vertical'"})
		return
	}

	upDownSites, err := u.GetStatusPageLastUpDownCount(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	html, err := getUpDownSitesiFrame(upDownSites, theme, apiURL, slug, containersDisplay, showTitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Data(http.StatusOK, "text/html", html)
}

func getUpDownSitesiFrame(upDownSites *UpDownSites, theme, apiURL, slug, containersDisplay string, showTitle bool) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="COLOR-SCHEME">
    <title>UptimeKuma iFrame</title>
    <style>
      ::-webkit-scrollbar {
        width: 7px;
      }

      ::-webkit-scrollbar-thumb {
        background-color: SCROLLBAR-THUMB-BACKGROUND-COLOR;
        border-radius: 2.3px;
      }

      ::-webkit-scrollbar-track {
        background-color: transparent;
      }

      ::-webkit-scrollbar-track:hover {
        background-color: SCROLLBAR-TRACK-BACKGROUND-COLOR;
      }
    </style>
    <style>
        body {
            background: transparent !important;
            margin: 0;
            padding: 0;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
        }

        .title-container {
            font-size: 30px;
            color: TITLE-COLOR;
            align-items: center;
            text-align: center;
            font-weight: bold;
            margin-bottom: 30px;
            width: 100%;
        }

        .info-containers {
            display: CONTAINERS-DISPLAY;
        }

        .info-container {
            box-sizing: border-box;
            border: 1px solid transparent;
            border-radius: 5px;
            background-color: rgba(9, 12, 16, 0.3);

            padding: 10px;
            margin: CONTAINER-MARGIN;
            width: 100%;

            align-items: center;
            justify-content: center;
            text-align: center;
            font-weight: bold;
        }

        .info-container:last-child {
            margin-bottom: 0 !important;
            margin-right: 0 !important;
        }

        .stats {
            color: STATS-COLOR;
        }


    </style>

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = 'API-URL/v1/hash/uptimekuma?slug=SLUG';
                const response = await fetch(url);
                const data = await response.json();

                if (lastHash === null) {
                    lastHash = data.hash;
                } else {
                    if (data.hash !== lastHash) {
                        lastHash = data.hash;
                        location.reload();
                    }
                }
            } catch (error) {
                console.error('Error getting last update from the API:', error);
            }
        }

        function fetchAndUpdate() {
            fetchData();
            setTimeout(fetchAndUpdate, 5000); // 5 seconds
        }

        fetchAndUpdate();
        
    </script>

</head>
<body>
<div>
    IFRAME-TITLE

    <div class="info-containers">
        <div class="info-container">
            <div class="stats">
                {{ .Up }}
            </div>
            <div>
                Up
            </div>
        </div>
        <div class="info-container">
            <div class="stats">
                {{ .Down }}
            </div>
            <div>
                Down
            </div>
        </div>
        <div class="info-container">
            <div class="stats">
                UPTIME-PERCENTAGE%
            </div>
            <div>
                Uptime
            </div>
        </div>
    </div>

</div>
</body>
</html>
    `
	// Homarr theme
	scrollbarThumbBackgroundColor := "rgba(209, 219, 227, 1)"
	scrollbarTrackBackgroundColor := "#ffffff"
	titleColor := "#000000"
	statsColor := "#5b6762"
	if theme == "dark" {
		scrollbarThumbBackgroundColor = "#484d64"
		scrollbarTrackBackgroundColor = "rgba(37, 40, 53, 1)"
		titleColor = "white"
		statsColor = "#949f9b"
	}

	if apiURL != "" {
		html = strings.ReplaceAll(html, "API-URL", apiURL)
		html = strings.ReplaceAll(html, "SLUG", slug)
	} else {
		html = strings.ReplaceAll(html, "fetchAndUpdate();", "// fetchAndUpdate")
	}

	if !showTitle {
		html = strings.ReplaceAll(html, "IFRAME-TITLE", "")
	} else {
		html = strings.ReplaceAll(html, "IFRAME-TITLE", `<div><div class="title-container">Uptime Kuma</div></div>`)
	}

	if containersDisplay == "horizontal" {
		html = strings.ReplaceAll(html, "CONTAINERS-DISPLAY", "flex")
		html = strings.ReplaceAll(html, "CONTAINER-MARGIN", "0px 10px 0px 0px")
	} else {
		html = strings.ReplaceAll(html, "CONTAINERS-DISPLAY", "block")
		html = strings.ReplaceAll(html, "CONTAINER-MARGIN", "0px 0px 10px 0px")
	}

	var uptimePercentage int
	if upDownSites.Up+upDownSites.Down != 0 {
		uptimePercentage = (upDownSites.Up * 100) / (upDownSites.Up + upDownSites.Down)
	} else {
		uptimePercentage = 0
	}
	html = strings.Replace(html, "COLOR-SCHEME", theme, -1)
	html = strings.Replace(html, "TITLE-COLOR", titleColor, -1)
	html = strings.Replace(html, "STATS-COLOR", statsColor, -1)
	html = strings.Replace(html, "UPTIME-PERCENTAGE", strconv.Itoa(uptimePercentage), -1)
	html = strings.Replace(html, "SCROLLBAR-THUMB-BACKGROUND-COLOR", scrollbarThumbBackgroundColor, -1)
	html = strings.Replace(html, "SCROLLBAR-TRACK-BACKGROUND-COLOR", scrollbarTrackBackgroundColor, -1)

	tmpl := template.Must(template.New("uptime").Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, upDownSites)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

// GetHash returns the hash of the up/down sites
func (u *UptimeKuma) GetHash(c *gin.Context) {
	slug := c.Query("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "slug must be provided"})
		return
	}

	upDownSites, err := u.GetStatusPageLastUpDownCount(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	hash := sources.GetHash(upDownSites, time.Now().Format("2006-01-02"))

	c.JSON(http.StatusOK, gin.H{"hash": fmt.Sprintf("%x", hash)})
}
