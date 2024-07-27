package netdata

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

var n *Netdata

var backgroundImageURL = "https://avatars.githubusercontent.com/u/43390781"

type Netdata struct {
	Address string
	Token   string
}

func New(address, token string) (*Netdata, error) {
	if n != nil {
		return n, nil
	}

	newN := &Netdata{}
	err := newN.Init(address, token)
	if err != nil {
		return nil, err
	}

	n = newN

	return n, nil
}

// Init sets the Netdata properties from the configs
func (n *Netdata) Init(address, token string) error {
	if address == "" || token == "" {
		return fmt.Errorf("NETDATA_ADDRESS and NETDATA_TOKEN variables should be set")
	}
	n.Address = strings.TrimSuffix(address, "/")
	n.Token = token

	return nil
}

// GetiFrame returns an HTML/CSS code to be used as an iFrame
func (n *Netdata) GetiFrame(c *gin.Context) {
	var err error
	theme := c.Query("theme")
	if theme == "" {
		theme = "light"
	} else if theme != "dark" && theme != "light" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "theme must be 'dark' or 'light'"})
		return
	}

	apiURL := c.Query("api_url")
	if apiURL != "" {
		_, err = url.ParseRequestURI(apiURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "api_url must be a valid URL like 'http://192.168.1.46:8080' or 'https://sub.domain.com'"})
			return
		}
	}

	queryLimit := c.Query("limit")
	var limit int
	if queryLimit == "" {
		limit = -1
	} else {
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "limit must be a number"})
			return
		}
	}

	alarms, err := n.GetAlarms(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var html []byte
	if len(alarms) < 1 {
		var apiURLPath string
		if apiURL != "" {
			apiURLPath = apiURL + "/v1/hash/netdata?limit=" + strconv.Itoa(limit)
		}
		html = sources.GetBaseNothingToShowiFrame("#226fff", backgroundImageURL, "center", "cover", "0.3", apiURLPath)
	} else {
		html, err = getAlarmsiFrame(n.Address, alarms, theme, apiURL, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Errorf("Couldn't create HTML code: %s", err.Error()))
			return
		}
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func getAlarmsiFrame(netdataAddress string, alarms []*Alarm, theme, apiURL string, limit int) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="{{ .Theme }}">
    <script src="https://kit.fontawesome.com/3f763b063a.js" crossorigin="anonymous"></script>
    <title>Netdata iFrame</title>
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
        body {
            background: transparent !important;
            margin: 0;
            padding: 0;
            width: calc(100% - 3px);
        }

        .alarms-container {
            height: 84px;

            position: relative;
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 14px;

            border-radius: 10px;
            border: 1px solid rgba(56, 58, 64, 1);
        }

        .background-image { 
            background-position: 50% 50%;
            background-size: 80%;
            background-image: url('https://avatars.githubusercontent.com/u/43390781');
            position: absolute;
            filter: brightness(0.3);
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            z-index: -1;
            border-radius: 10px;
        }

        .text-wrap {
            flex-grow: 1;
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
            width: 1px !important;
            margin-right: 10px;
            padding: 20px;

            /* this set the ellipsis (...) properties only if the attributes below are overwritten*/
            color: white;
            font-weight: bold;
        }

        .alarm-summary {
            font-size: 15px;
            color: white;
            font-family: -apple-system, BtaskMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
            text-decoration: none;
        }

        .alarm-summary:hover {
            text-decoration: underline;
        }

        .more-info-container {
            flex-grow: 1;
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
            margin-right: 10px;
            margin-top: 3px;

            /* this set the ellipsis (...) properties only if the attributes below are overwritten*/
            color: #99b6bb;
            font-weight: bold;
        }

        .info-label {
            text-decoration: none;
            font-family: ui-sans-serif, system-ui, -apple-system, BtaskMacSystemFont,
              Segoe UI, Roboto, Helvetica Neue, Arial, Noto Sans, sans-serif, Apple Color Emoji,
              Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji;
            font-feature-settings: normal;
            font-variation-settings: normal;
            font-weight: 600;
            font-size: 1rem;
            line-height: 1.5rem;

            margin-right: 7px;
        }

        a.info-label:hover {
            text-decoration: underline;
        }

        .alarm-info-container {
            display: flex;
            flex-direction: column;
            padding: 20px;
            justify-content: center;
            align-items: center;
            min-width: 91.33px;
        }

        .alarm-value-label {
            text-decoration: none;
            font-family: ui-sans-serif, system-ui, -apple-system, BtaskMacSystemFont,
              Segoe UI, Roboto, Helvetica Neue, Arial, Noto Sans, sans-serif, Apple Color Emoji,
              Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji;
            font-feature-settings: normal;
            font-variation-settings: normal;
            font-weight: 600;
            color: white;
            font-size: 1rem;
            line-height: 1.5rem;
            margin: 0 0 5px 0;
        }

        .alarm-status-label {
            text-decoration: none;
            font-family: ui-sans-serif, system-ui, -apple-system, BtaskMacSystemFont,
              Segoe UI, Roboto, Helvetica Neue, Arial, Noto Sans, sans-serif, Apple Color Emoji,
              Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji;
            font-feature-settings: normal;
            font-variation-settings: normal;
            font-weight: 700;
            font-size: 0.6875rem;
            line-height: calc(1.125rem);

            padding: 0px calc(0.666667rem) 0px calc(0.666667rem) !important;

            display: inline-block;
            border-radius: 1rem;
            margin: 0;
        }
    </style>

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = '{{ .APIURL }}/v1/hash/netdata?limit={{ .APILimit }}';
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
            setTimeout(fetchAndUpdate, 10000); // 10 seconds
        }

        {{ if .APIURL }}
            fetchAndUpdate();
        {{ end }}
    </script>

</head>
<body>
{{ range .Alarms }}
    <div class="alarms-container">
        <div class="background-image"></div>

        <div class="text-wrap">
            <i class="fa-solid fa-bell"></i> <a href="{{ with . }}{{ $.NetdataAddress }}{{ end }}" target="_blank" class="alarm-summary">{{ .Summary }}</a>
            <div class="more-info-container">
                <span class="info-label"><i class="fa-solid fa-calendar-days"></i> {{ .LastStatusChange.Format "2006-01-02 15h04" }}</span> 
                <span class="info-label"><i class="fa-solid fa-gear"></i> {{ .Type }}</span>
                <span class="info-label"><i class="fa-solid fa-cube"></i> {{ .Component }}</span>
            </div>
        </div>
    
        <div class="alarm-info-container">
            <p class="alarm-value-label">{{ .ValueString }}</p>
            <div>
                <p class="alarm-status-label" style="color: white; background-color: {{ getStatusColor .Status }};">{{ .Status }}</p>
            </div>
        </div>
    </div>
{{ end }}
</body>
</html>
	`
	// Homarr theme
	scrollbarThumbBackgroundColor := "#d1dbe3"
	scrollbarTrackBackgroundColor := "#ffffff"
	if theme == "dark" {
		scrollbarThumbBackgroundColor = "#484d64"
		scrollbarTrackBackgroundColor = "rgba(37, 40, 53, 1)"
	}

	templateData := iframeTemplateData{
		Alarms:                        alarms,
		Theme:                         theme,
		APIURL:                        apiURL,
		APILimit:                      limit,
		BackgroundImageURL:            backgroundImageURL,
		NetdataAddress:                netdataAddress,
		ScrollbarThumbBackgroundColor: scrollbarThumbBackgroundColor,
		ScrollbarTrackBackgroundColor: scrollbarTrackBackgroundColor,
	}

	templateFuncs := template.FuncMap{
		"getStatusColor": func(status string) string {
			switch status {
			case "CLEAR":
				return "green"
			case "WARNING":
				return "orange"
			case "CRITICAL":
				return "red"
			default:
				return "gray"
			}
		},
	}

	tmpl := template.Must(template.New("alarms").Funcs(templateFuncs).Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, &templateData)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

type iframeTemplateData struct {
	Alarms                        []*Alarm
	Theme                         string
	APIURL                        string
	APILimit                      int
	BackgroundImageURL            string
	NetdataAddress                string
	ScrollbarThumbBackgroundColor string
	ScrollbarTrackBackgroundColor string
}

// GetHash returns the hash of the alarms
func (n *Netdata) GetHash(c *gin.Context) {
	queryLimit := c.Query("limit")
	var limit int
	var err error
	if queryLimit == "" {
		limit = -1
	} else {
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "limit must be a number"})
			return
		}
	}

	alarms, err := n.GetAlarms(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	hash := sources.GetHash(alarms, time.Now().Format("2006-01-02"))

	c.JSON(http.StatusOK, gin.H{"hash": fmt.Sprintf("%x", hash)})
}
