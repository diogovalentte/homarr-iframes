package linkwarden

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources"
)

var backgroundImageURL = "https://avatars.githubusercontent.com/u/135248736?s=280&v=4"

type Linkwarden struct {
	Address string
	Token   string
}

func (l *Linkwarden) Init() error {
	address := config.GlobalConfigs.LinkwardenConfigs.Address
	token := config.GlobalConfigs.LinkwardenConfigs.Token
	if address == "" || token == "" {
		return fmt.Errorf("LINKWARDEN_ADDRESS and LINKWARDEN_TOKEN variables should be set")
	}

	if strings.HasSuffix(address, "/") {
		address = address[:len(address)-1]
	}

	l.Address = address
	l.Token = token

	return nil
}

func (l *Linkwarden) GetiFrame(c *gin.Context) {
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
		}
	}

	collectionID := c.Query("collectionId")
	var url string
	if collectionID != "" {
		url = l.Address + "/api/v1/links?collectionId=" + collectionID
	} else {
		url = l.Address + "/api/v1/links"
	}

	links := map[string][]*Link{}
	err = baseRequest(url, &links)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("Error while doing API request: %s", err.Error())})
		return
	}

	res, exists := links["response"]
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("No 'response' field in API response")})
		return
	}

	if limit >= 0 {
		res = res[:limit]
	}

	var html []byte
	if len(res) < 1 {
		html = sources.GetBaseNothingToShowiFrame(theme, backgroundImageURL, "center", "cover", "0.3")
	} else {
		html, err = getLinksiFrame(l.Address, res, theme)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("Couldn't create HTML code: %s", err.Error())})
			return
		}
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func getLinksiFrame(linkwardenAddress string, links []*Link, theme string) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="LINKS-CONTAINER-BACKGROUND-COLOR">
    <script src="https://kit.fontawesome.com/3f763b063a.js" crossorigin="anonymous"></script>
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
            width: calc(100% - 3px);
        }

        .links-container {
            height: 84px;

            position: relative;
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 14px;

            border-radius: 10px;
            border: 1px solid rgba(56, 58, 64, 1);
        }

        .links-container img {
            padding: 20px;
        }

        .background-image {
            background-position: 50% 47.2%;
            background-size: cover;
            position: absolute;
            filter: brightness(0.3);
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            z-index: -1;
            border-radius: 10px;
        }

        .link-icon {
            width: 32px;
            height: 32px;
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

        .link-name {
            font-size: 15px;
            color: white;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
            text-decoration: none;
        }

        .link-name:hover {
            text-decoration: underline;
        }

        .info-label {
            text-decoration: none;
            font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont,
              Segoe UI, Roboto, Helvetica Neue, Arial, Noto Sans, sans-serif, Apple Color Emoji,
              Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji;
            font-feature-settings: normal;
            font-variation-settings: normal;
            font-weight: 600;
            color: #4f6164;
            font-size: 1rem;
            line-height: 1.5rem;
        }

        a.info-label:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
{{ range . }}
    <div class="links-container">

        <div style="background-image: url('LINKS-CONTAINER-BACKGROUND-IMAGE');" class="background-image"></div>

        <img class="link-icon" src="https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url={{ .URL }}/&size=32" alt="Link Site Favicon">

        <div class="text-wrap">
            {{ if .Name }}
                <a href="{{ .URL }}" target="_blank" class="link-name">{{ .Name }}</a>
            {{ else if .Description }}
                <a href="{{ .URL }}" target="_blank" class="link-name">{{ .Description }}</a>
            {{ else }}
                <a href="{{ .URL }}" target="_blank" class="link-name">&lt;No name or description&gt;</a>
            {{ end }}

            <div>
                <span style="margin-right: 7px;" class="info-label"><i class="fa-solid fa-calendar-days"></i> {{ .CreatedAt.Format "Jan 2, 2006" }}</span>
                {{ if .CollectionID }}
                    <a href="LINKWARDEN-ADDRESS/collections/{{ .CollectionID }}" target="_blank" class="info-label"><i style="color: {{ .Collection.Color }};" class="fa-solid fa-folder-closed"></i> {{ .Collection.Name }}</a>
                {{ end }}
            </div>
        </div>

    </div>
{{ end }}
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

	html = strings.Replace(html, "LINKWARDEN-ADDRESS", linkwardenAddress, -1)
	html = strings.Replace(html, "LINKS-CONTAINER-BACKGROUND-COLOR", theme, -1)
	html = strings.Replace(html, "LINKS-CONTAINER-BACKGROUND-IMAGE", backgroundImageURL, -1)
	html = strings.Replace(html, "SCROLLBAR-THUMB-BACKGROUND-COLOR", scrollbarThumbBackgroundColor, -1)
	html = strings.Replace(html, "SCROLLBAR-TRACK-BACKGROUND-COLOR", scrollbarTrackBackgroundColor, -1)

	tmpl := template.Must(template.New("links").Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, links)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
