package media

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources"
)

func getIframeData(radarrReleaseType string, unmonitored bool) (*Calendar, error) {
	var isAnySourceValid bool
	calendar := &Calendar{}
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 1)

	if config.GlobalConfigs.Radarr.Address != "" && config.GlobalConfigs.Radarr.APIKey != "" {
		isAnySourceValid = true
		radarr, err := NewRadarr(config.GlobalConfigs.Radarr.Address, config.GlobalConfigs.Radarr.InternalAddress, config.GlobalConfigs.Radarr.APIKey)
		if err != nil {
			return nil, fmt.Errorf("couldn't create Radarr client: %s", err.Error())
		}
		radarrCalendar, err := radarr.GetCalendar(unmonitored, startDate, endDate, radarrReleaseType)
		if err != nil {
			return nil, fmt.Errorf("couldn't get Radarr calendar: %s", err.Error())
		}
		for _, release := range radarrCalendar.Releases {
			calendar.Releases = append(calendar.Releases, release)
		}
	}

	if config.GlobalConfigs.Sonarr.Address != "" && config.GlobalConfigs.Sonarr.APIKey != "" {
		isAnySourceValid = true
		sonarr, err := NewSonarr(config.GlobalConfigs.Sonarr.Address, config.GlobalConfigs.Sonarr.InternalAddress, config.GlobalConfigs.Sonarr.APIKey)
		if err != nil {
			return nil, fmt.Errorf("couldn't create Sonarr client: %s", err.Error())
		}
		sonarrCalendar, err := sonarr.GetCalendar(unmonitored, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("couldn't get Sonarr calendar: %s", err.Error())
		}
		for _, release := range sonarrCalendar.Releases {
			calendar.Releases = append(calendar.Releases, release)
		}
	}

	if !isAnySourceValid {
		return nil, fmt.Errorf("no valid source found. Please check the docs for what environment variables should be set")
	}

	return calendar, nil
}

// GetiFrame returns an HTML/CSS code to be used as an iFrame
func GetiFrame(c *gin.Context) {
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

	var (
		radarrReleaseType string
		showUnmonitored   bool
	)
	showEpisodeHours := true

	radarrReleaseType = c.Query("radarrReleaseType")
	if radarrReleaseType == "" {
		radarrReleaseType = "inCinemas"
	}

	queryShowUnmonitored := c.Query("showUnmonitored")
	switch queryShowUnmonitored {
	case "true":
		showUnmonitored = true
	case "false":
	case "":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "showUnmonitored must be empty, 'true', or 'false'"})
		return
	}

	queryShowEpisodeHours := c.Query("showEpisodesHour")
	switch queryShowEpisodeHours {
	case "true":
	case "false":
		showEpisodeHours = false
	case "":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "showEpisodeHours must be empty, 'true', or 'false'"})
		return
	}

	iframeRequestData, err := getIframeData(radarrReleaseType, showUnmonitored)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var html []byte
	html, err = getMediaReleasesiFrame(iframeRequestData, theme, apiURL, radarrReleaseType, showUnmonitored, showEpisodeHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("couldn't create iFrame: %s", err.Error()))
		return
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func getMediaReleasesiFrame(calendar *Calendar, theme string, apiURL string, radarrReleaseType string, showUnmonitored, showEpisodeHours bool) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="{{ .Theme }}">
    <script src="https://kit.fontawesome.com/3f763b063a.js" crossorigin="anonymous"></script>
    <title>Media Releases iFrame</title>
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

        .releases-container {
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
            background-position: 50% 49.5%;
            background-size: 105%;
            position: absolute;
            filter: brightness(0.3);
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            z-index: -1;
            border-radius: 10px;
        }

        .release-cover {
            border-radius: 2px;
            object-fit: cover;
            width: 30px;
            height: 50px;
        }

        img.release-cover {
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

        .release-title {
            font-size: 15px;
            color: white;
            font-family: -apple-system, BtaskMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
            text-decoration: none;
        }

        .release-title:hover {
            text-decoration: underline;
        }

        .more-info-container {
            flex-grow: 1;
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
            margin-right: 10px;

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

        .source-info-container {
            display: flex;
            flex-direction: column;
            padding: 20px;
            justify-content: center;
            align-items: center;
            min-width: 91.33px;
        }

        .source-label {
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
            margin: 0 0 5px 0;
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
                var url = '{{ .APIURL }}/v1/hash/media_releases?showUnmonitored={{ .APIShowUnmonitored }}&radarrReleaseType={{ .APIRadarrReleaseType }}';
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
{{ range .Calendar.Releases }}
    <div class="releases-container">

        <div class="background-image" style="background-image: url('{{ .CoverImageURL }}');"></div>
        <img
            class="release-cover"
            src="{{ .CoverImageURL }}"
            alt="Media Release Poster"
        />

        <div class="text-wrap">
            {{ if eq .Source "Sonarr" }}
                <a href="{{ with . }}{{ $.SonarrAddress }}{{ end }}/series/{{ .Slug }}" target="_blank" class="release-title">{{ .Title }}</a>
                <div class="more-info-container">
                    {{ with . }}{{ if $.ShowEpisodeHours }}
                        <span class="info-label" style="display: inline-block; min-width: 63.25px;"><i class="fa-solid fa-calendar-days"></i> {{ .ReleaseDate.Format "15h04" }}</span>
                    {{ end }}{{ end }}
                    <span class="info-label"><i class="fas fa-tv fa-xm"></i> S{{ .EpisodeDetails.SeasonNumber }}E{{ .EpisodeDetails.EpisodeNumber}} - {{ .EpisodeDetails.EpisodeName }}</span>
                </div>
            {{ else if eq .Source "Radarr" }}
                <a href="{{ with . }}{{ $.RadarrAddress }}{{ end }}/movie/{{ .Slug }}" target="_blank" class="release-title">{{ .Title }}</a>
            {{ end }}
        </div>

        <div class="source-info-container">
            <p class="source-label" style="color: {{ getSourceColor .Source }};">{{ .Source }}</p>
            <div>
                {{ if .IsDownloaded }}
                    <p class="status-label" style="color: white; background-color: green;">Available</p>
                {{ else }}
                    {{ if .ShouldBeDownloaded }}
                        <p class="status-label" style="color: white; background-color: red;">Not Available</p>
                    {{ else }}
                        <p class="status-label" style="color: white; background-color: #99b6bb;">Not Available</p>
                    {{ end }}
                {{ end }}
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
		Calendar:                      calendar,
		Theme:                         theme,
		APIURL:                        apiURL,
		APIShowUnmonitored:            showUnmonitored,
		APIRadarrReleaseType:          radarrReleaseType,
		ShowEpisodeHours:              showEpisodeHours,
		SonarrAddress:                 strings.TrimSuffix(config.GlobalConfigs.Sonarr.Address, "/"),
		RadarrAddress:                 strings.TrimSuffix(config.GlobalConfigs.Radarr.Address, "/"),
		ScrollbarThumbBackgroundColor: scrollbarThumbBackgroundColor,
		ScrollbarTrackBackgroundColor: scrollbarTrackBackgroundColor,
	}

	templateFuncs := template.FuncMap{
		"getSourceColor": func(source string) string {
			switch source {
			case "Sonarr":
				return "#1c7ed6"
			case "Radarr":
				return "#f59f00"
			default:
				return "#99b6bb"
			}
		},
	}

	tmpl, err := template.New("releases").Funcs(templateFuncs).Parse(html)
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, &templateData)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

type iframeTemplateData struct {
	Calendar                      *Calendar
	Theme                         string
	APIURL                        string
	APIShowUnmonitored            bool
	APIRadarrReleaseType          string
	ShowEpisodeHours              bool
	SonarrAddress                 string
	RadarrAddress                 string
	ScrollbarThumbBackgroundColor string
	ScrollbarTrackBackgroundColor string
}

// GetHash returns the hash of the media releases
func GetHash(c *gin.Context) {
	radarrReleaseType := c.Query("radarrReleaseType")
	if radarrReleaseType == "" {
		radarrReleaseType = "inCinemas"
	}
	queryShowUnmonitored := c.Query("showUnmonitored")
	var showUnmonitored bool
	switch queryShowUnmonitored {
	case "true":
		showUnmonitored = true
	case "false":
	case "":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "showUnmonitored must be empty, 'true', or 'false'"})
		return
	}

	releases, err := getIframeData(radarrReleaseType, showUnmonitored)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	hash := sources.GetHash(*releases, time.Now().Format("2006-01-02"))

	c.JSON(http.StatusOK, gin.H{"hash": fmt.Sprintf("%x", hash)})
}
