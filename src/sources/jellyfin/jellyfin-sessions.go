package jellyfin

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/diogovalentte/homarr-iframes/src/sources"
	"github.com/gin-gonic/gin"
)

func (j *Jellyfin) GetSessionsiFrame(c *gin.Context) {
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

	activeSince := c.Query("activeWithinSeconds")
	var activeWithinSeconds int
	if activeSince == "" {
		activeWithinSeconds = 0
	} else {
		activeWithinSeconds, err = strconv.Atoi(activeSince)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "activeWithinSeconds must be a number"})
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

	sessions, err := j.GetSessions(limit, activeWithinSeconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("couldn't get items: %s", err.Error())})
		return
	}

	var html []byte
	if len(sessions) < 1 {
		var apiURLPath string
		if apiURL != "" {
			apiURLPath = apiURL + "/v1/hash/jellyfin/sessions?limit=" + strconv.Itoa(limit) + "&theme=" + theme
			if activeWithinSeconds > 0 {
				apiURLPath += "&activeWithinSeconds=" + strconv.Itoa(activeWithinSeconds)
			}
		}
		backgroundImgURL := "https://avatars.githubusercontent.com/u/45698031?s=280&v=4"
		html = sources.GetBaseNothingToShowiFrame(theme, backgroundImgURL, "center", "cover", "brightness(0.3)", apiURLPath)
	} else {
		html, err = j.getSessionsiFrame(sessions, theme, apiURL, limit, activeWithinSeconds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("couldn't create HTML code: %s", err.Error())})
			return
		}
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func (j *Jellyfin) getSessionsiFrame(sessions []*Session, theme, apiURL string, limit, activeWithinSeconds int) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="{{ .Theme }}">
    <script src="https://kit.fontawesome.com/3f763b063a.js" crossorigin="anonymous"></script>
    <title>Jellyfin Sessions</title>
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

        .episode-info {
            font-size: 0.85em;
            color: #99b6bb;
            font-weight: normal;
            margin-left: 5px;
            font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont,
              Segoe UI, Roboto, Helvetica Neue, Arial, Noto Sans, sans-serif, Apple Color Emoji,
              Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji;
            font-feature-settings: normal;
            font-variation-settings: normal;
            text-decoration: none;
        }

        .episode-info:hover {
            text-decoration: underline;
        }
        
        /* User avatar styles */
        .user-avatar {
            object-fit: cover;
            width: 25px;
            height: 25px;
            border-radius: 50%;
            margin: 0 10px;
        }
        
        .user-container {
            display: inline-flex;
            align-items: center;
            margin-right: 20px;
        }
        
        .username {
            text-decoration: none;
            font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont,
              Segoe UI, Roboto, Helvetica Neue, Arial, Noto Sans, sans-serif, Apple Color Emoji,
              Segoe UI Emoji, Segoe UI Symbol, Noto Color Emoji;
            font-feature-settings: normal;
            font-variation-settings: normal;
            font-weight: 600;
            color: #99b6bb;
            font-size: 0.9rem;
        }
    </style>

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = '{{ .APIURL }}/v1/hash/jellyfin/sessions?limit={{ .APILimit }}{{ if .ActiveWithinSeconds }}&activeWithinSeconds={{ .ActiveWithinSeconds }}{{ end }}';
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
            setTimeout(fetchAndUpdate, 3000); // 3 seconds
        }

        {{ if .APIURL }}
            fetchAndUpdate();
        {{ end }}
    </script>

</head>
<body>
{{ range .Sessions }}
    {{ if .NowPlayingItem }}
    <div class="items-container">
        <div class="background-image" style="background-image: url('{{ .NowPlayingItem.BackdropImageURL }}');"></div>
        <img
            class="item-cover"
            src="{{ .NowPlayingItem.PrimaryImageURL }}"
            alt="Media Item Cover"
        />

        <div class="text-wrap">
            {{ if and (eq .NowPlayingItem.Type "Episode") .NowPlayingItem.SeriesName }}
                <a href="{{ .NowPlayingItem.ItemURL }}" target="_blank" class="item-title" title="{{ .NowPlayingItem.SeriesName }}">
                    {{ .NowPlayingItem.SeriesName }}
                </a>
                {{ if and .NowPlayingItem.SeasonNumber .NowPlayingItem.EpisodeNumber }}
                    <a href="{{ .NowPlayingItem.EpisodeURL }}" target="_blank" class="episode-info" title="Go to episode">
                        S{{ printf "%02d" .NowPlayingItem.SeasonNumber }}E{{ printf "%02d" .NowPlayingItem.EpisodeNumber }}
                    </a>
                {{ end }}
            {{ else }}
                <a href="{{ .NowPlayingItem.ItemURL }}" target="_blank" class="item-title" title="{{ .NowPlayingItem.Name }}">{{ .NowPlayingItem.Name }}</a>
            {{ end }}
            
            <div class="labels-div">
                {{ if and .PlayState.PositionTicks .NowPlayingItem.RunTimeTicks }}
                    <span class="info-label">
                        {{ if .PlayState.IsPaused }}
                            <i class="fa-solid fa-pause"></i>
                        {{ else }}
                            <i class="fa-solid fa-play"></i>
                        {{ end }}
                        {{ formatTime .PlayState.PositionTicks }}/{{ formatTime .NowPlayingItem.RunTimeTicks }}
                    </span>
                {{ end }}
            </div>
        </div>
        
        <div class="user-container">
            <img 
                class="user-avatar"
                src="{{ .UserAvatarURL }}"
                alt="{{ .UserName }} Avatar"
                onerror="this.onerror=null; this.src='https://avatars.githubusercontent.com/u/45698031?s=280&v=4';"
            />
            <span class="username">{{ .UserName }}</span>
        </div>
    </div>
    {{ end }}
{{ end }}

`
	// Homarr theme
	scrollbarThumbBackgroundColor := "#d1dbe3"
	scrollbarTrackBackgroundColor := "#ffffff"
	if theme == "dark" {
		scrollbarThumbBackgroundColor = "#484d64"
		scrollbarTrackBackgroundColor = "rgba(37, 40, 53, 1)"
	}

	templateData := sessionsiframeTemplateData{
		Sessions:                      sessions,
		Theme:                         theme,
		APIURL:                        apiURL,
		APILimit:                      limit,
		ActiveWithinSeconds:           activeWithinSeconds,
		ScrollbarThumbBackgroundColor: scrollbarThumbBackgroundColor,
		ScrollbarTrackBackgroundColor: scrollbarTrackBackgroundColor,
	}

	funcMap := template.FuncMap{
		"formatTime": func(ticks int64) string {
			totalSeconds := ticks / 10000000

			hours := totalSeconds / 3600
			minutes := (totalSeconds % 3600) / 60
			seconds := totalSeconds % 60

			return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		},
	}

	tmpl := template.Must(template.New("sessions").Funcs(funcMap).Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, &templateData)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

type sessionsiframeTemplateData struct {
	Theme                         string
	Sessions                      []*Session
	APIURL                        string
	ScrollbarThumbBackgroundColor string
	ScrollbarTrackBackgroundColor string
	APILimit                      int
	ActiveWithinSeconds           int
}

func (j *Jellyfin) GetSessionsHash(c *gin.Context) {
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

	activeWithinSecondsStr := c.Query("activeWithinSeconds")
	var activeWithinSeconds int
	if activeWithinSecondsStr == "" {
		activeWithinSeconds = 0
	} else {
		activeWithinSeconds, err = strconv.Atoi(activeWithinSecondsStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "activeWithinSeconds must be a number"})
			return
		}
	}

	sessions, err := j.GetSessions(limit, activeWithinSeconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	hash := sources.GetHash(sessions, time.Now().Format("2006-01-02"))

	c.JSON(http.StatusOK, gin.H{"hash": fmt.Sprintf("%x", hash)})
}
