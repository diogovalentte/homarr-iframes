package uptimekuma

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

// UptimeKuma is the UptimeKuma source
type UptimeKuma struct {
	// UptimeKuma instance URL
	Address string
}

// Init sets the UptimeKuma properties from the configs
func (u *UptimeKuma) Init() error {
	address := config.GlobalConfigs.UptimeKumaConfigs.Address
	if address == "" {
		return fmt.Errorf("UPTIMEKUMA_ADDRESS variable should be set")
	}

	if strings.HasSuffix(address, "/") {
		address = address[:len(address)-1]
	}

	u.Address = address

	return nil
}

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

	upDownSites, err := u.GetStatusPageLastUpDownCount(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	html, err := getUpDownSitesiFrame(upDownSites, theme, apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Data(http.StatusOK, "text/html", html)
}

func getUpDownSitesiFrame(upDownSites *UpDownSites, theme, apiURL string) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="COLOR-SCHEME">
    <title>Movie Display Template</title>
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
        }

        .title-container {
            font-size: 30px;
            color: white;
            font-family: -apple-system, BtaskMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
            align-items: center;
            text-align: center;
            font-weight: bold;
        }

        .info-containers {
            display: flex;
        }

        .info-container {
            box-sizing: border-box;
            border: 1px solid red;

            padding: 10px;
            margin: 5px;
            width: 33.33%;

            align-items: center;
            justify-content: center;
            text-align: center;
        }


    </style>

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = 'API-URL/v1/hash/uptimekuma';
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
    <div>
        <div class="title-container" style="margin-bottom: 20px;">Uptime Monitor</div>
    </div>

    <div class="info-containers">
        <div class="info-container" style="color: green">
            <div>
                {{ .Up }}
            </div>
            <div>
                Ola
            </div>
        </div>
        <div class="info-container" style="color: red">{{ .Down }}</div>
        <div class="info-container" style="color: blue">UPTIME-PERCENTAGE</div>
    </div>

</div>
</body>
</html>
    `
	// Homarr theme
	scrollbarThumbBackgroundColor := "rgba(209, 219, 227, 1)"
	scrollbarTrackBackgroundColor := "#ffffff"
	if theme == "dark" {
		scrollbarThumbBackgroundColor = "#484d64"
		scrollbarTrackBackgroundColor = "rgba(37, 40, 53, 1)"
	}

	if apiURL != "" {
		html = strings.Replace(html, "API-URL", apiURL, -1)
	} else {
		html = strings.Replace(html, "fetchAndUpdate();", "// fetchAndUpdate", -1)
	}

	upDownSites.Down = 5
	uptimePercentage := (upDownSites.Up * 100) / (upDownSites.Up + upDownSites.Down)

	html = strings.Replace(html, "COLOR-SCHEME", theme, -1)
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
