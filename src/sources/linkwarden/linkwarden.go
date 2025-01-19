package linkwarden

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

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources"
)

var (
	defaultBackgroundImgURL = "https://avatars.githubusercontent.com/u/135248736?s=280&v=4"
	l                       *Linkwarden
)

type Linkwarden struct {
	Address          string
	InternalAddress  string
	Token            string
	BackgroundImgURL string
}

func New() (*Linkwarden, error) {
	if l != nil {
		return l, nil
	}

	address := config.GlobalConfigs.Linkwarden.Address
	internalAddress := config.GlobalConfigs.Linkwarden.InternalAddress
	token := config.GlobalConfigs.Linkwarden.Token
	backgroundImgURL := config.GlobalConfigs.Vikunja.BackgroundImgURL
	if backgroundImgURL == "" {
		backgroundImgURL = defaultBackgroundImgURL
	}

	newL := &Linkwarden{}
	err := newL.Init(address, internalAddress, token, backgroundImgURL)
	if err != nil {
		return nil, err
	}

	l = newL

	return l, nil
}

func (l *Linkwarden) Init(address, internalAddress, token, backgroundImgURL string) error {
	if address == "" || token == "" {
		return fmt.Errorf("LINKWARDEN_ADDRESS and LINKWARDEN_TOKEN variables should be set")
	}

	l.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		l.InternalAddress = l.Address
	} else {
		l.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	l.Token = token
	l.BackgroundImgURL = backgroundImgURL

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
			return
		}
	}

	collectionID := c.Query("collectionId")

	backgroundPosition := c.Query("background_position")
	if backgroundPosition == "" {
		backgroundPosition = "50% 47.2%"
	}
	backgroundSize := c.Query("background_size")
	if backgroundSize == "" {
		backgroundSize = "cover"
	}
	backgroundFilter := c.Query("background_filter")
	if backgroundFilter == "" {
		backgroundFilter = "brightness(0.3)"
	}

	apiURL := c.Query("api_url")
	if apiURL != "" {
		_, err = url.ParseRequestURI(apiURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "api_url must be a valid URL like 'http://192.168.1.46:8080' or 'https://sub.domain.com'"})
			return
		}
	}

	links, err := l.GetLinks(limit, collectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("couldn't get links: %s", err.Error())})
		return
	}

	var html []byte
	if len(links) < 1 {
		var apiURLPath string
		if apiURL != "" {
			apiURLPath = apiURL + "/v1/hash/linkwarden?limit=" + strconv.Itoa(limit) + "&collectionId=" + collectionID
		}
		html = sources.GetBaseNothingToShowiFrame(theme, l.BackgroundImgURL, "center", "cover", backgroundFilter, apiURLPath)
	} else {
		html, err = l.getLinksiFrame(links, theme, l.BackgroundImgURL, backgroundPosition, backgroundSize, backgroundFilter, apiURL, collectionID, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("couldn't create HTML code: %s", err.Error())})
			return
		}
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func (l *Linkwarden) getLinksiFrame(links []*Link, theme, backgroundImgURL, backgroundPosition, backgroundSize, backgroundFilter, apiURL, collectionID string, limit int) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="{{ .Theme }}">
    <script src="https://kit.fontawesome.com/3f763b063a.js" crossorigin="anonymous"></script>
    <title>Linkwarden iFrame</title>
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

        .links-container {
            height: 84px;

            position: relative;
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin: 8.50px;

            border-radius: 10px;
            border: 1px solid rgba(56, 58, 64, 1);
        }

        .links-container img {
            padding: 20px;
        }

        .background-image {
            background-image: url('{{ .BackgroundImageURL }}');
            background-position: {{ .BackgroundPosition }};
            background-size: {{ .BackgroundSize }};
            filter: {{ .BackgroundFilter }};
            position: absolute;
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

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = '{{ .APIURL }}/v1/hash/linkwarden?limit={{ .APILimit }}&collectionId={{ .CollectionID }}';
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
{{ range .Links }}
    <div class="links-container">

        <div class="background-image"></div>

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
                    <i style="color: {{ .Collection.Color }};" class="fa-solid fa-folder-closed"></i> <a href="{{ with . }}{{ $.LinkwardenAddress }}{{ end }}/collections/{{ .CollectionID }}" target="_blank" class="info-label">{{ .Collection.Name }}</a>
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

	templateData := iframeTemplateData{
		Links:                         links,
		Theme:                         theme,
		APIURL:                        apiURL,
		APILimit:                      limit,
		LinkwardenAddress:             l.Address,
		BackgroundImageURL:            backgroundImgURL,
		BackgroundPosition:            template.CSS(backgroundPosition),
		BackgroundSize:                template.CSS(backgroundSize),
		BackgroundFilter:              template.CSS(backgroundFilter),
		ScrollbarThumbBackgroundColor: scrollbarThumbBackgroundColor,
		ScrollbarTrackBackgroundColor: scrollbarTrackBackgroundColor,
		CollectionID:                  collectionID,
	}

	tmpl := template.Must(template.New("links").Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, templateData)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

type iframeTemplateData struct {
	Theme                         string
	APIURL                        string
	LinkwardenAddress             string
	BackgroundImageURL            string
	BackgroundPosition            template.CSS
	BackgroundSize                template.CSS
	BackgroundFilter              template.CSS
	ScrollbarThumbBackgroundColor string
	ScrollbarTrackBackgroundColor string
	CollectionID                  string
	Links                         []*Link
	APILimit                      int
}

// GetHash returns the hash of the bookmarks
func (l *Linkwarden) GetHash(c *gin.Context) {
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

	collectionID := c.Query("collectionId")

	pLinks, err := l.GetLinks(limit, collectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("couldn't get links: %s", err.Error())})
		return
	}

	var links []Link
	for _, link := range pLinks {
		link.Description = nil
		link.CollectionID = nil
		link.Collection = nil
		links = append(links, *link)
	}

	hash := sources.GetHash(links, time.Now().Format("2006-01-02"))

	c.JSON(http.StatusOK, gin.H{"hash": fmt.Sprintf("%x", hash)})
}
