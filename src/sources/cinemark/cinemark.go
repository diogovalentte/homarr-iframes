package cinemark

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/diogovalentte/homarr-iframes/src/sources"
)

var backgroundImageURL = "https://static.vecteezy.com/system/resources/previews/025/470/292/large_2x/background-image-date-at-the-cinema-popcorn-ai-generated-photo.jpeg"

type Cinemark struct{}

// GetiFrame returns an iframe with the in theater movies for a specific city
func (_ *Cinemark) GetiFrame(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city is required"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a number"})
		}
	}

	theaters := c.Query("theaters")

	theme := c.Query("theme")
	if theme == "" {
		theme = "light"
	} else if theme != "dark" && theme != "light" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "theme must be 'dark' or 'light'"})
		return
	}

	scraper := Scraper{}
	movies, err := scraper.GetInTheatersMovies(city, limit, theaters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var html []byte
	if len(movies) < 1 {
		html = sources.GetBaseNothingToShowiFrame(theme, backgroundImageURL, "center", "cover", "0.3")
	} else {
		html, err = getMoviesiFrame(movies, theme)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Errorf("Couldn't create HTML code: %s", err.Error()))
			return
		}
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func getMoviesiFrame(movies []Movie, theme string) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
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
            background-color: MOVIES-CONTAINER-BACKGROUND-COLOR;
            margin: 0;
            padding: 0;
        }

        .movies-container {
            width: calc(100% - MOVIES-CONTAINER-WIDTHpx);
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
            background-image: url('MOVIES-CONTAINER-BACKGROUND-IMAGE');
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
            width: 42.08px;
            height: 42.08px;
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
</head>
<body>
{{ range . }}
        <div class="movies-container">

            <div class="background-image"></div>

            <img
                class="movie-cover"
                src="{{ .CoverImgURL }}"
                alt="Movie Cover"
            />

            <div class="text-wrap">
                <a href="{{ .URL }}" target="_blank" class="movie-name">{{.Title}}</a>
            </div>

            <div>
                {{ if .Label }}
                    <div class="movie-label">{{ .Label }}</div>
                {{end}}

                <div style="background-color: {{ .AgeRatingColor }}" class="movie-age-rating">{{ .AgeRating }}</div>
            </div>

        </div>

    </div>
{{ end }}
</body>
</html>
    `
	// Set the container width based on the number of movies for better fitting with Homarr
	containerWidth := "1.6"
	if len(movies) > 3 {
		containerWidth = "8"
	}

	// Homarr theme
	containerBackgroundColor := "#ffffff"
	scrollbarThumbBackgroundColor := "rgba(209, 219, 227, 1)"
	scrollbarTrackBackgroundColor := "#ffffff"
	if theme == "dark" {
		containerBackgroundColor = "#25262b"
		scrollbarThumbBackgroundColor = "#484d64"
		scrollbarTrackBackgroundColor = "rgba(37, 40, 53, 1)"
	}

	html = strings.Replace(html, "MOVIES-CONTAINER-WIDTH", containerWidth, -1)
	html = strings.Replace(html, "MOVIES-CONTAINER-BACKGROUND-COLOR", containerBackgroundColor, -1)
	html = strings.Replace(html, "MOVIES-CONTAINER-BACKGROUND-IMAGE", backgroundImageURL, -1)
	html = strings.Replace(html, "SCROLLBAR-THUMB-BACKGROUND-COLOR", scrollbarThumbBackgroundColor, -1)
	html = strings.Replace(html, "SCROLLBAR-TRACK-BACKGROUND-COLOR", scrollbarTrackBackgroundColor, -1)

	tmpl := template.Must(template.New("movies").Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, movies)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
