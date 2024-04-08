package vikunja

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

var backgroundImageURL = "https://vikunja.io/images/vikunja.png"

type Vikunja struct {
	Address  string
	Username string
	Password string
}

func (v *Vikunja) Init() error {
	address := config.GlobalConfigs.VikunjaConfigs.Address
	username := config.GlobalConfigs.VikunjaConfigs.Username
	password := config.GlobalConfigs.VikunjaConfigs.Password
	if address == "" || username == "" || password == "" {
		return fmt.Errorf("VIKUNJA_ADDRESS, VIKUNJA_USERNAME, and VIKUNJA_PASSWORD variables should be set")
	}

	if strings.HasSuffix(address, "/") {
		address = address[:len(address)-1]
	}

	v.Address = address
	v.Username = username
	v.Password = password

	return nil
}

func (v *Vikunja) GetiFrame(c *gin.Context) {
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

	tasks := []*Task{}
	if limit != 0 {
		tasks, err = v.GetTasks(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	var html []byte
	if len(tasks) < 1 {
		html = sources.GetBaseNothingToShowiFrame("#226fff", backgroundImageURL, "center", "cover", "0.3")
	} else {
		html, err = getTasksiFrame(v.Address, tasks, theme)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Errorf("Couldn't create HTML code: %s", err.Error()))
			return
		}
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func getTasksiFrame(vikunjaAddress string, tasks []*Task, theme string) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
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
            background-color: TASKS-CONTAINER-BACKGROUND-COLOR;
            margin: 0;
            padding: 0;
        }

        .tasks-container {
            width: calc(100% - TASKS-CONTAINER-WIDTHpx);
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
            background-image: url('TASKS-CONTAINER-BACKGROUND-IMAGE');
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

        .task-title {
            font-size: 15px;
            color: white;
            font-family: -apple-system, BtaskMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
            text-decoration: none;
        }

        .task-title:hover {
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
    </style>
</head>
<body>
{{ range . }}
    <div class="tasks-container">

        <div class="background-image"></div>

        <div style="padding-left: 20px;" class="text-wrap">
            <a href="VIKUNJA-ADDRESS/tasks/{{ .ID }}" target="_blank" class="task-title">{{ .Title }}</a>

            <div>

                <span class="info-label"><i class="fa-solid fa-calendar-days"></i> Created: {{ .CreatedAt.Format "Jan 2, 2006" }}</span>

                {{ if not .DueDate.IsZero }}
                    <span class="info-label"><i class="fa-solid fa-calendar-days"></i> Due: {{ .DueDate.Format "Jan 2, 2006" }}</span>
                {{ else if not .EndDate.IsZero }}
                    <span class="info-label"><i class="fa-solid fa-calendar-days"></i> End: {{ .EndDate.Format "Jan 2, 2006" }}</span>
                {{ end }}

                {{ if eq .Priority 3 }}
                    <span style="color: #ff851b;" class="info-label">! High</span>
                {{ else if eq .Priority 4 }}
                    <span style="color: #ff4136;" class="info-label">! Urgent</span>
                {{ else if eq .Priority 5 }}
                    <span style="color: #ff4136;" class="info-label">! DO NOW !</span>
                {{ end }}

            </div>

        </div>

    </div>
{{ end }}
</body>
</html>
	`
	// Set the container width based on the number of tasks for better fitting with Homarr
	containerWidth := "1.6"
	if len(tasks) > 3 {
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

	html = strings.Replace(html, "VIKUNJA-ADDRESS", vikunjaAddress, -1)
	html = strings.Replace(html, "TASKS-CONTAINER-WIDTH", containerWidth, -1)
	html = strings.Replace(html, "TASKS-CONTAINER-BACKGROUND-COLOR", containerBackgroundColor, -1)
	html = strings.Replace(html, "TASKS-CONTAINER-BACKGROUND-IMAGE", backgroundImageURL, -1)
	html = strings.Replace(html, "SCROLLBAR-THUMB-BACKGROUND-COLOR", scrollbarThumbBackgroundColor, -1)
	html = strings.Replace(html, "SCROLLBAR-TRACK-BACKGROUND-COLOR", scrollbarTrackBackgroundColor, -1)

	tmpl := template.Must(template.New("tasks").Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, tasks)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
