package jellyfin

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/diogovalentte/homarr-iframes/src/sources"
	"github.com/gin-gonic/gin"
)

func (j *Jellyfin) GetRecentlyAddediFrame(c *gin.Context) {
	theme := c.Query("theme")
	if theme == "" {
		theme = "light"
	} else if theme != "dark" && theme != "light" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "theme must be 'dark' or 'light'"})
		return
	}

	showLimit := c.Query("limit")
	var limit int
	var err error
	if showLimit == "" {
		limit = 0
	} else {
		limit, err = strconv.Atoi(showLimit)
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

	userId := c.Query("userId")
	if userId != "" {
		match, err := regexp.MatchString(`^[0-9a-fA-F]{32}$`, userId)
		if err != nil || !match {
			c.JSON(http.StatusBadRequest, gin.H{"message": "userId must be a valid Jellyfin user ID (32 hexadecimal characters)"})
			return
		}
	}

	parentId := c.Query("parentId")
	if parentId != "" {
		match, err := regexp.MatchString(`^[0-9a-fA-F]{32}$`, parentId)
		if err != nil || !match {
			c.JSON(http.StatusBadRequest, gin.H{"message": "parentId must be a valid Jellyfin folder/library ID (32 hexadecimal characters)"})
			return
		}
	}

	includeItemTypes := c.Query("includeItemTypes")
	if includeItemTypes != "" {
		allowedTypes := map[string]bool{
			"AggregateFolder": true, "Audio": true, "AudioBook": true, "BasePluginFolder": true,
			"Book": true, "BoxSet": true, "Channel": true, "ChannelFolderItem": true,
			"CollectionFolder": true, "Episode": true, "Folder": true, "Genre": true,
			"ManualPlaylistsFolder": true, "Movie": true, "LiveTvChannel": true, "LiveTvProgram": true,
			"MusicAlbum": true, "MusicArtist": true, "MusicGenre": true, "MusicVideo": true,
			"Person": true, "Photo": true, "PhotoAlbum": true, "Playlist": true,
			"PlaylistsFolder": true, "Program": true, "Recording": true, "Season": true,
			"Series": true, "Studio": true, "Trailer": true, "TvChannel": true,
			"TvProgram": true, "UserRootFolder": true, "UserView": true, "Video": true, "Year": true,
		}
		itemTypesList := strings.Split(includeItemTypes, ",")
		for _, itemType := range itemTypesList {
			itemType = strings.TrimSpace(itemType)
			if !allowedTypes[itemType] {
				allowedTypesList := make([]string, 0, len(allowedTypes))
				for t := range allowedTypes {
					allowedTypesList = append(allowedTypesList, t)
				}
				sort.Strings(allowedTypesList)

				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("'%s' is not a valid item type. Allowed types are: %s",
						itemType, strings.Join(allowedTypesList, ", ")),
				})
				return
			}
		}
	}

	jellyfinQueryLimit := c.Query("queryLimit")
	var queryLimit int
	if jellyfinQueryLimit == "" {
		queryLimit = 0
	} else {
		queryLimit, err = strconv.Atoi(jellyfinQueryLimit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "queryLimit must be a number"})
			return
		}
	}

	items, err := j.GetLatestItems(limit, queryLimit, userId, parentId, includeItemTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("couldn't get items: %s", err.Error())})
		return
	}

	var html []byte
	if len(items) < 1 {
		var apiURLPath string
		if apiURL != "" {
			apiURLPath = apiURL + "/v1/hash/jellyfin?limit=" + strconv.Itoa(limit) + "&theme=" + theme
			if userId != "" {
				apiURLPath += "&userId=" + userId
			}
			if parentId != "" {
				apiURLPath += "&parentId=" + parentId
			}
			if includeItemTypes != "" {
				apiURLPath += "&includeItemTypes=" + includeItemTypes
			}
			if queryLimit > 0 {
				apiURLPath += "&queryLimit=" + strconv.Itoa(queryLimit)
			}
		}
		backgroundImgURL := "https://avatars.githubusercontent.com/u/45698031?s=280&v=4"
		html = sources.GetBaseNothingToShowiFrame(theme, backgroundImgURL, "center", "cover", "brightness(0.3)", apiURLPath)
	} else {
		html, err = j.getItemsiFrame(items, theme, apiURL, limit, userId, parentId, includeItemTypes, queryLimit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("couldn't create HTML code: %s", err.Error())})
			return
		}
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func (j *Jellyfin) getItemsiFrame(items []*Item, theme, apiURL string, limit int, userId, parentId, includeItemTypes string, queryLimit int) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="{{ .Theme }}">
    <script src="https://kit.fontawesome.com/3f763b063a.js" crossorigin="anonymous"></script>
    <title>Jellyfin Recently Added</title>
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

        .items-container {
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
            background-position: 50% 15%;
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

        .item-cover {
            border-radius: 2px;
            object-fit: cover;
            width: 30px;
            height: 50px;
            padding: 20px;
        }

        .text-wrap {
            flex-grow: 1;
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
            width: 1px !important;
            margin-right: 10px;
            color: white;
            font-weight: bold;
        }

        .item-title {
            font-size: 15px;
            color: white;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
            text-decoration: none;
        }

        .item-title:hover {
            text-decoration: underline;
        }

        .labels-div {
            min-height: 24px;
            display: flex;
            align-items: center;
        }

        .info-label {
            text-decoration: none;
            font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont,
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

        .type-label {
            text-decoration: none;
            font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont,
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
    </style>

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = '{{ .APIURL }}/v1/hash/jellyfin/recently?limit={{ .APILimit }}{{ if .UserId }}&userId={{ .UserId }}{{ end }}{{ if .ParentId }}&parentId={{ .ParentId }}{{ end }}{{ if .IncludeItemTypes }}&includeItemTypes={{ .IncludeItemTypes }}{{ end }}{{ if .QueryLimit }}&queryLimit={{ .QueryLimit }}{{ end }}';
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
{{ range .Items }}
    <div class="items-container">
        <div class="background-image" style="background-image: url('{{ .BackdropImageURL }}');"></div>
        <img
            class="item-cover"
            src="{{ .PrimaryImageURL }}"
            alt="Media Item Cover"
        />

        <div class="text-wrap">
            <a href="{{ .ItemURL }}" target="_blank" class="item-title" title="{{ .Name }}">{{ .Name }}</a>
            <div class="labels-div">
                {{ if .Year }}
                    <span class="info-label"><i class="fa-solid fa-calendar-days"></i> {{ .Year }}</span>
                {{ end }}
            </div>
        </div>
    </div>
{{ end }}

`
	// Homarr theme
	scrollbarThumbBackgroundColor := "#d1dbe3"
	scrollbarTrackBackgroundColor := "#ffffff"
	if theme == "dark" {
		scrollbarThumbBackgroundColor = "#484d64"
		scrollbarTrackBackgroundColor = "rgba(37, 40, 53, 1)"
	}

	templateData := iframeTemplateData{
		Items:                         items,
		Theme:                         theme,
		APIURL:                        apiURL,
		APILimit:                      limit,
		UserId:                        userId,
		ParentId:                      parentId,
		IncludeItemTypes:              includeItemTypes,
		QueryLimit:                    queryLimit,
		ScrollbarThumbBackgroundColor: scrollbarThumbBackgroundColor,
		ScrollbarTrackBackgroundColor: scrollbarTrackBackgroundColor,
	}

	tmpl := template.Must(template.New("items").Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, &templateData)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

type iframeTemplateData struct {
	Theme                         string
	Items                         []*Item
	APIURL                        string
	ScrollbarThumbBackgroundColor string
	ScrollbarTrackBackgroundColor string
	UserId                        string
	ParentId                      string
	IncludeItemTypes              string
	APILimit                      int
	QueryLimit                    int
}

func (j *Jellyfin) GetHash(c *gin.Context) {
	frameQueryLimit := c.Query("limit")
	var limit int
	var err error
	if frameQueryLimit == "" {
		limit = -1
	} else {
		limit, err = strconv.Atoi(frameQueryLimit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "limit must be a number"})
			return
		}
	}

	userId := c.Query("userId")
	parentId := c.Query("parentId")
	includeItemTypes := c.Query("includeItemTypes")

	jellyfinQueryLimit := c.Query("queryLimit")
	var queryLimit int
	if jellyfinQueryLimit == "" {
		queryLimit = 0
	} else {
		queryLimit, err = strconv.Atoi(jellyfinQueryLimit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "queryLimit must be a number"})
			return
		}
	}

	items, err := j.GetLatestItems(limit, queryLimit, userId, parentId, includeItemTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	hash := sources.GetHash(items, time.Now().Format("2006-01-02"))

	c.JSON(http.StatusOK, gin.H{"hash": fmt.Sprintf("%x", hash)})
}
