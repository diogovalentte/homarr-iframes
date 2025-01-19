package mediarequets

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
	"github.com/diogovalentte/homarr-iframes/src/sources/jellyseerr"
	"github.com/diogovalentte/homarr-iframes/src/sources/overseerr"
)

// GetiFrame returns an HTML/CSS code to be used as an iFrame
func GetiFrame(c *gin.Context) {
	theme := c.Query("theme")
	if theme == "" {
		theme = "light"
	} else if theme != "dark" && theme != "light" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "theme must be 'dark' or 'light'"})
		return
	}

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

	apiURL := c.Query("api_url")
	if apiURL != "" {
		_, err = url.ParseRequestURI(apiURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "api_url must be a valid URL like 'http://192.168.1.46:8080' or 'https://sub.domain.com'"})
			return
		}
	}

	var (
		filter                string
		sort                  string
		requestedByOverseerr  int
		requestedByJellyseerr int
	)
	filter = c.Query("filter")
	sort = c.Query("sort")
	requestedByOverseerrStr := c.Query("requestedByOverseerr")
	if requestedByOverseerrStr != "" {
		requestedByOverseerr, err = strconv.Atoi(requestedByOverseerrStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "requestedByOverseerr must be a number"})
			return
		}
	}

	requestedByJeellyseerrStr := c.Query("requestedByJellyseerr")
	if requestedByJeellyseerrStr != "" {
		requestedByJellyseerr, err = strconv.Atoi(requestedByJeellyseerrStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "requestedByJellyseerr must be a number"})
			return
		}
	}

	iframeRequestData, err := getIframeData(limit, filter, sort, requestedByOverseerr, requestedByJellyseerr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var html []byte
	html, err = getRequestsiFrame(iframeRequestData, theme, apiURL, limit, filter, sort, requestedByOverseerr, requestedByJellyseerr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("couldn't create HTML code: %s", err.Error()))
		return
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func getRequestsiFrame(requests []overseerr.IframeRequestData, theme, apiURL string, limit int, filter, sort string, requestedByOverseerr, requestedByJellyseerr int) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="{{ .Theme }}">
    <script src="https://kit.fontawesome.com/3f763b063a.js" crossorigin="anonymous"></script>
    <title>Overseerr iFrame</title>
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
        }

        .requests-container {
            height: 84px;

            position: relative;
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin: 8.50px;

            border-radius: 10px;
            border: 1px solid rgba(56, 58, 64, 1);
        }

        .background-image { 
            background-position: 50% 49.5%;
            background-size: 100%;
            position: absolute;
            filter: brightness(0.3);
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            z-index: -1;
            border-radius: 10px;
        }

        .request-cover {
            border-radius: 2px;
            object-fit: cover;
            width: 30px;
            height: 50px;
        }

        img.request-cover {
            padding: 20px;
        }

        .text-wrap {
            flex-grow: 1;
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
            width: 1px !important;
            margin-right: 10px;

            /* this set the ellipsis (...) properties only if the attributes below are overwritten*/
            color: white;
            font-weight: bold;
        }

        .request-title {
            font-size: 15px;
            color: white;
            font-family: -apple-system, BtaskMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
            text-decoration: none;
        }

        .request-title:hover {
            text-decoration: underline;
        }

        .labels-div {
            min-height: 24px;
            display: flex;
            align-items: center;
        }

        .info-label {
            text-decoration: none;
            font-family: ui-sans-serif, system-ui, -apple-system, BtaskMacSystemFont,
              Segoe UI, Roboto, Helvetica Neue, Arial, Noto Sans, sans-serif, Apple Color Emoji,
              Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji;
            font-feature-settings: normal;
            font-variation-settings: normal;
            font-weight: 600;
            color: #99b6bb;
            font-size: 1rem;
            line-height: 1.5rem;

            margin-right: 7px;
        }

        a.info-label:hover {
            text-decoration: underline;
        }

        .status-label {
            text-decoration: none;
            font-family: ui-sans-serif, system-ui, -apple-system, BtaskMacSystemFont,
              Segoe UI, Roboto, Helvetica Neue, Arial, Noto Sans, sans-serif, Apple Color Emoji,
              Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji;
            font-feature-settings: normal;
            font-variation-settings: normal;
            font-weight: 700;
            font-size: 0.6875rem;
            line-height: calc(1.125rem);
            text-transform: uppercase;

            padding: 0px calc(0.666667rem) 0px calc(0.666667rem) !important;

            display:inline-block;
            border-radius: 1rem;
            padding: 0.1rem 0.5rem;
        }

        .requested-by-container {
            display: inline-block;
            text-align: center;
            margin: 20px 20px 20px 10px;
        }

        .requested-by-avatar {
            object-fit: cover;
            width: 25px;
            height: 25px;
            border-radius: 50%;
        }

        .username {
            text-decoration: none;
            font-family: ui-sans-serif, system-ui, -apple-system, BtaskMacSystemFont,
              Segoe UI, Roboto, Helvetica Neue, Arial, Noto Sans, sans-serif, Apple Color Emoji,
              Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji;
            font-feature-settings: normal;
            font-variation-settings: normal;
            font-weight: 600;
            color: #99b6bb;
            font-size: 1rem;
            line-height: 1.5rem;
        }

        a.username:hover {
            text-decoration: underline;
        }
    </style>

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = '{{ .APIURL }}/v1/hash/media_requests?limit={{ .APILimit }}&filter={{ .APIFilter }}&sort={{ .APISort }}&requestedByOverseerr={{ .APIRequestedByOverseerr }}&requestedByJellyseerr={{ .APIRequestedByJellyseerr }}';
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
{{ range .Requests }}
    <div class="requests-container">

        <div class="background-image" style="background-image: url('{{ .Media.BackdropURL }}');"></div>
        <img
            class="request-cover"
            src="{{ .Media.PosterURL }}"
            alt="Media Request Cover"
        />

        <div class="text-wrap">
            <a href="{{ .Media.URL }}" target="_blank" class="request-title">{{ .Media.Name }}</a>
            <div class="labels-div">
                {{ if .Media.Year }}
                    <span class="info-label"><i class="fa-solid fa-calendar-days"></i> {{ .Media.Year }}</span>
                {{ end }}
                <span class="status-label" style="color: {{ .Status.Color }}; background-color: {{ .Status.BackgroundColor }} ">{{ .Status.Status }}</span>
            </div>
        </div>

        <img
            class="requested-by-avatar"
            src="{{ .Request.AvatarURL }}"
            alt="Requested By Avatar"
        />
        <div class="requested-by-container">
            <a href="{{ .Request.UserProfileURL }}" target="_blank" class="username">{{ .Request.Username }}</a>
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
		Requests:                      requests,
		Theme:                         theme,
		APIURL:                        apiURL,
		APILimit:                      limit,
		APIFilter:                     filter,
		APISort:                       sort,
		APIRequestedByOverseerr:       requestedByOverseerr,
		APIRequestedByJellyseerr:      requestedByJellyseerr,
		ScrollbarThumbBackgroundColor: scrollbarThumbBackgroundColor,
		ScrollbarTrackBackgroundColor: scrollbarTrackBackgroundColor,
	}

	tmpl := template.Must(template.New("requests").Parse(html))

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
	APIFilter                     string
	APISort                       string
	ScrollbarThumbBackgroundColor string
	ScrollbarTrackBackgroundColor string
	Requests                      []overseerr.IframeRequestData
	APIRequestedByOverseerr       int
	APIRequestedByJellyseerr      int
	APILimit                      int
}

// GetHash returns the hash of the requests
func GetHash(c *gin.Context) {
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

	var (
		filter                string
		sort                  string
		requestedByOverseerr  int
		requestedByJellyseerr int
	)
	filter = c.Query("filter")
	sort = c.Query("sort")
	requestedByOverseerrStr := c.Query("requestedByOverseerr")
	if requestedByOverseerrStr != "" {
		requestedByOverseerr, err = strconv.Atoi(requestedByOverseerrStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "requestedByOverseerr must be a number"})
			return
		}
	}

	requestedByJeellyseerrStr := c.Query("requestedByJellyseerr")
	if requestedByJeellyseerrStr != "" {
		requestedByJellyseerr, err = strconv.Atoi(requestedByJeellyseerrStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "requestedByJellyseerr must be a number"})
			return
		}
	}

	iframeRequestData, err := getIframeData(limit, filter, sort, requestedByOverseerr, requestedByJellyseerr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	hash := sources.GetHash(iframeRequestData, time.Now().Format("2006-01-02"))

	c.JSON(http.StatusOK, gin.H{"hash": fmt.Sprintf("%x", hash)})
}

func getIframeData(limit int, filter, sort string, requestedByOverseerr, requestedByJellyseerr int) ([]overseerr.IframeRequestData, error) {
	var requests []overseerr.IframeRequestData

	o, err := overseerr.New()
	if err != nil {
		if !strings.Contains(err.Error(), "variables should be set") {
			return nil, err
		}
	} else {
		ORequests, err := o.GetIframeData(limit, filter, sort, requestedByOverseerr)
		if err != nil {
			return nil, err
		}
		requests = append(requests, ORequests...)
	}
	j, err := jellyseerr.New()
	if err != nil {
		if !strings.Contains(err.Error(), "variables should be set") {
			return nil, err
		}
	} else {
		JRequests, err := j.GetIframeData(limit, filter, sort, requestedByJellyseerr)
		if err != nil {
			return nil, err
		}
		requests = append(requests, JRequests...)
	}

	return requests, nil
}
