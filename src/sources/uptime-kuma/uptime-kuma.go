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
    <meta name="color-scheme" content="{{ .Theme }}">
    <title>UptimeKuma iFrame</title>
    <style>
      ::-webkit-scrollbar {
        width: 7px;
      }

      ::-webkit-scrollbar-thumb {
        background-color: {{ .ScrollbarThumbBackgroundColor }};
        border-radius: 2.3px;
      }

      ::-webkit-scrollbar-track {
        background-color: transparent;
      }

      ::-webkit-scrollbar-track:hover {
        background-color: {{ .ScrollbarTrackBackgroundColor }};
      }
    </style>
    <style>
        {{ .CSSCode }}
    </style>

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = '{{ .APIURL }}/v1/hash/uptimekuma?slug={{ .Slug }}';
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

        {{ if .APIURL }}
            fetchAndUpdate();
        {{ end }}
    </script>

</head>
<body>
<div class="main">
    {{ .Title }}

    <div class="info-containers">
        <div class="info-container">
            <div class="stats">
                {{ .UpSites }}
            </div>
            <div>
                Up
            </div>
        </div>
        <div class="info-container">
            <div class="stats">
                {{ .DownSites }}
            </div>
            <div>
                Down
            </div>
        </div>
        <div class="info-container">
            <div class="stats">
                {{ .UptimePercentage }}%
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

	var CSSCode string
	if containersDisplay == "horizontal" {
		if showTitle {
			CSSCode = `
                body {
                    background: transparent !important;
                    margin: 0;
                    padding: 0;
                    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    width: 100vw;
                    overflow: hidden;
                }

                div.main {
                    display: flex;
                    flex-direction: column;
                    justify-content: center;
                    align-items: center;
                    width: 100%;
                    height: 100%;
                    max-width: 1200px;
                    padding: 0;
                    box-sizing: border-box;
                    text-align: center;
                    margin-bottom: 20px;
                }

                .title-container {
                    display: flex;
                    font-size: 30px;
                    color: {{ .TitleColor }};
                    align-items: center;
                    justify-content: center;
                    text-align: center;
                    font-weight: bold;
                    width: 100%;
                    height: 50%;
                }

                .info-containers {
                    display: flex;
                    justify-content: space-between;
                    width: 100%;
                    height: 50%;
                    box-sizing: border-box;
                }

                .info-container {
                    display: flex;
                    flex-direction: column;
                    height: 100%;
                    flex: 1;
                    align-items: center;
                    justify-content: center;
                    text-align: center;
                    font-weight: bold;
                    padding: 10px;
                    margin: 5px;
                    box-sizing: border-box;
                    border: 1px solid transparent;
                    border-radius: 5px;
                    background-color: rgba(9, 12, 16, 0.3);
                }

                .stats {
                    color: {{ .StatsColor }};
                }
            `
		} else {
			CSSCode = `
                body {
                    background: transparent !important;
                    margin: 0;
                    padding: 0;
                    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    width: 100vw;
                    overflow: hidden;
                }

                div.main {
                    display: flex;
                    flex-direction: column;
                    justify-content: center;
                    align-items: center;
                    width: 100%;
                    height: 100%;
                    max-width: 1200px;
                    padding: 0;
                    box-sizing: border-box;
                    text-align: center;
                }

                .info-containers {
                    display: flex;
                    justify-content: space-between;
                    width: 100%;
                    box-sizing: border-box;
                }

                .info-container {
                    flex: 1;
                    align-items: center;
                    justify-content: center;
                    text-align: center;
                    font-weight: bold;
                    padding: 10px;
                    margin: 5px;
                    box-sizing: border-box;
                    border: 1px solid transparent;
                    border-radius: 5px;
                    background-color: rgba(9, 12, 16, 0.3);
                }

                .stats {
                    color: {{ .StatsColor }};
                }
            `
		}
	} else {
		if showTitle {
			CSSCode = `
                body {
                    background: transparent !important;
                    margin: 0;
                    padding: 0;
                    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    width: 100vw;
                    overflow: hidden;
                }

                div.main {
                    height: 100%;
                    width: 100%;
                }

                .title-container {
                    font-size: 30px;
                    color: {{ .TitleColor }};
                    text-align: center;
                    font-weight: bold;
                    margin: 0;
                    padding: 0;
                    height: 25%;
                    line-height: normal;
                    width: 100%;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                }

                .info-containers {
                    display: flex;
                    flex-direction: column;
                    width: 100%;
                    height: 75%;
                    box-sizing: border-box;
                }

                .info-container {
                    height: 33.33%;
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                    justify-content: center;
                    text-align: center;
                    font-weight: bold;
                    padding: 10px;
                    margin: 5px 0;
                    box-sizing: border-box;
                    border: 1px solid transparent;
                    border-radius: 5px;
                    background-color: rgba(9, 12, 16, 0.3);
                }

                .stats {
                    color: {{ .StatsColor }};
                }
            `
		} else {
			CSSCode = `
                body {
                    background: transparent !important;
                    margin: 0;
                    padding: 0;
                    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    width: 100vw;
                    overflow: hidden;
                }

                div.main {
                    height: 100%;
                    width: 100%;
                }

                .info-containers {
                    display: flex;
                    flex-direction: column;
                    width: 100%;
                    height: 100%;
                    box-sizing: border-box;
                }

                .info-container {
                    height: 33.33%;
                    display: flex;
                    flex-direction: column;

                    align-items: center;
                    justify-content: center;
                    text-align: center;
                    font-weight: bold;
                    padding: 10px;
                    margin: 5px 0;
                    box-sizing: border-box;
                    border: 1px solid transparent;
                    border-radius: 5px;
                    background-color: rgba(9, 12, 16, 0.3);
                }

                .stats {
                    color: {{ .StatsColor }};
                }
            `
		}
	}

	var uptimePercentage int
	if upDownSites.Up+upDownSites.Down != 0 {
		uptimePercentage = (upDownSites.Up * 100) / (upDownSites.Up + upDownSites.Down)
	} else {
		uptimePercentage = 0
	}

	templateData := iframeTemplateData{
		Theme:                         theme,
		APIURL:                        apiURL,
		Slug:                          slug,
		ScrollbarThumbBackgroundColor: scrollbarThumbBackgroundColor,
		ScrollbarTrackBackgroundColor: scrollbarTrackBackgroundColor,
		CSSCode:                       template.CSS(CSSCode),
		TitleColor:                    titleColor,
		Title:                         template.HTML(`<div class="title-container">Uptime Kuma</div>`),
		StatsColor:                    statsColor,
		UptimePercentage:              uptimePercentage,
		UpSites:                       upDownSites.Up,
		DownSites:                     upDownSites.Down,
	}

	if !showTitle {
		templateData.Title = ""
	}

	tmpl := template.Must(template.New("uptime").Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, &templateData)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

type iframeTemplateData struct {
	Theme                         string
	APIURL                        string
	Slug                          string
	ScrollbarThumbBackgroundColor string
	ScrollbarTrackBackgroundColor string
	CSSCode                       template.CSS
	TitleColor                    string
	Title                         template.HTML
	StatsColor                    string
	UptimePercentage              int
	UpSites                       int
	DownSites                     int
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
