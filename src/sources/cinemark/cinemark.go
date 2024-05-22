package cinemark

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

var backgroundImageURL = "https://static.vecteezy.com/system/resources/previews/025/470/292/large_2x/background-image-date-at-the-cinema-popcorn-ai-generated-photo.jpeg"

type Cinemark struct{}

// GetiFrame returns an iframe with the in theater movies for a specific city
func (_ *Cinemark) GetiFrame(c *gin.Context) {
	theaterIdsStr := c.Query("theaterIds")
	var theaterIds []int
	if theaterIdsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "theaterIds is required"})
		return
	} else {
		theaterStrings := strings.Split(theaterIdsStr, ",")
		for _, theaterStr := range theaterStrings {
			theaterId, err := strconv.Atoi(theaterStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "theaterIds must be a list of numbers"})
				return
			}
			theaterIds = append(theaterIds, theaterId)
		}
	}

	queryLimit := c.Query("limit")
	var limitProvided bool
	var limit int
	var err error
	if queryLimit != "" {
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "limit must be a number"})
			return
		}
		limitProvided = true
	}

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

	cinemark := Cinemark{}
	movies, err := cinemark.GetOnDisplayByTheater(theaterIds, limit, limitProvided)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var html []byte
	if len(movies) < 1 {
		var apiURLPath string
		if apiURL != "" {
			apiURLPath = apiURL + "/v1/hash/cinemark?limit=" + strconv.Itoa(limit) + "&theaterId=" + theaterIdsStr
		}
		html = sources.GetBaseNothingToShowiFrame(theme, backgroundImageURL, "center", "cover", "0.3", apiURLPath)
	} else {
		html, err = getMoviesiFrame(movies, theme, apiURL, limit, theaterIdsStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("Couldn't create HTML code: %s", err.Error())})
			return
		}
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func getMoviesiFrame(movies []Movie, theme, apiURL string, limit int, theaterIds string) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="{{ .Theme }}">
    <title>Movie Display Template</title>
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

        .movies-container {
            height: 84px;

            position: relative;
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 14px;

            border-radius: 10px;
            border: 1px solid rgba(56, 58, 64, 1);
        }

        .movies-container img {
            padding: 20px;
        }

        .background-image {
            background-image: url('{{ .BackgroundImageURL }}');
            background-position: center;
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

        .movie-cover {
            border-radius: 2px;
            object-fit: cover;
            width: 30px;
            height: 50px;
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

        .movie-name {
            font-size: 15px;
            color: white;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
            text-decoration: none;
        }

        .movie-name:hover {
            text-decoration: underline;
        }

        .more-info-container {
            display: flex;
            flex-direction: column;
            margin-left: auto;
            margin-right: 10px;
            width: 160px;
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

        .movie-label {
            display: inline-block;
            padding: 8px 10px;
            margin: 20px;
            background-color: rgb(150, 109, 109, 0.5);
            color: rgb(230, 101, 101);

            text-decoration: none; /* Remove underline */
            border-radius: 5px;
            font-size: 20px;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
              Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
        }

        .movie-age-rating {
            display: inline-block;
            padding: 7px 10px;
            margin-right: 20px;
            min-width: 42.08px;
            min-height: 42.08px;
            border-radius: 5px;
            box-sizing: border-box;

            font-size: 20px;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
              Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
            color: white;
            font-weight: 800;
            text-align: center;
        }
    </style>

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = '{{ .APIURL }}/v1/hash/cinemark?limit={{ .APILimit }}&theaterId={{ .TheaterId }}';
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
{{ range .Movies }}
        <div class="movies-container">

            <div class="background-image"></div>

            <img
                class="movie-cover"
                src="{{ .CoverImgURL }}"
                alt="Movie Cover"
            />

            <div class="text-wrap">
                <a href="{{ .URL }}" target="_blank" class="movie-name">{{ .Name }}</a>
            </div>

            <div>
                {{ if .Genre }}
                    <div class="movie-label">{{ .Genre }}</div>
                {{end}}

                {{ if .AgeRating }}
                    <div style="background-color: {{ .AgeRatingColor }}" class="movie-age-rating">{{ .AgeRating }}</div>
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

	// if apiURL != "" {
	// 	html = strings.Replace(html, "API-URL", apiURL, -1)
	// 	html = strings.Replace(html, "API-LIMIT", strconv.Itoa(limit), -1)
	// 	html = strings.Replace(html, "API-THEATER-ID", theaterId, -1)
	// } else {
	// 	html = strings.Replace(html, "fetchAndUpdate();", "// fetchAndUpdate", -1)
	// }
	//
	// html = strings.Replace(html, "MOVIES-CONTAINER-BACKGROUND-COLOR", theme, -1)
	// html = strings.Replace(html, "MOVIES-CONTAINER-BACKGROUND-IMAGE", backgroundImageURL, -1)
	// html = strings.Replace(html, "SCROLLBAR-THUMB-BACKGROUND-COLOR", scrollbarThumbBackgroundColor, -1)
	// html = strings.Replace(html, "SCROLLBAR-TRACK-BACKGROUND-COLOR", scrollbarTrackBackgroundColor, -1)

	templateData := iframeTemplateData{
		Movies:                        movies,
		Theme:                         theme,
		APIURL:                        apiURL,
		APILimit:                      limit,
		TheaterId:                     theaterIds,
		BackgroundImageURL:            backgroundImageURL,
		ScrollbarThumbBackgroundColor: scrollbarThumbBackgroundColor,
		ScrollbarTrackBackgroundColor: scrollbarTrackBackgroundColor,
	}

	tmpl := template.Must(template.New("movies").Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, &templateData)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

type iframeTemplateData struct {
	Movies                        []Movie
	Theme                         string
	APIURL                        string
	APILimit                      int
	TheaterId                     string
	BackgroundImageURL            string
	ScrollbarThumbBackgroundColor string
	ScrollbarTrackBackgroundColor string
}

// GetHash returns the hash of the in theater movies for a specific city
func (_ *Cinemark) GetHash(c *gin.Context) {
	theaterIdsStr := c.Query("theaterIds")
	var theaterIds []int
	if theaterIdsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "theaterIds is required"})
		return
	} else {
		theaterStrings := strings.Split(theaterIdsStr, ",")
		for _, theaterStr := range theaterStrings {
			theaterId, err := strconv.Atoi(theaterStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "theaterIds must be a list of numbers"})
				return
			}
			theaterIds = append(theaterIds, theaterId)
		}
	}

	queryLimit := c.Query("limit")
	var limitProvided bool
	var limit int
	var err error
	if queryLimit != "" {
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "limit must be a number"})
			return
		}
		limitProvided = true
	}

	cinemark := Cinemark{}
	movies, err := cinemark.GetOnDisplayByTheater(theaterIds, limit, limitProvided)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	hash := sources.GetHash(movies, time.Now().Format("2006-01-02"))

	c.JSON(http.StatusOK, gin.H{"hash": fmt.Sprintf("%x", hash)})
}
