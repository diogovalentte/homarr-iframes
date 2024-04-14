package vikunja

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

var backgroundImageURL = "https://vikunja.io/images/vikunja.png"

// Vikunja is the a source
type Vikunja struct {
	Address string
	Token   string
}

// Init sets the Vikunja properties from the configs
func (v *Vikunja) Init() error {
	address := config.GlobalConfigs.VikunjaConfigs.Address
	token := config.GlobalConfigs.VikunjaConfigs.Token
	if address == "" || token == "" {
		return fmt.Errorf("VIKUNJA_ADDRESS and VIKUNJA_TOKEN variables should be set")
	}

	if strings.HasSuffix(address, "/") {
		address = address[:len(address)-1]
	}

	v.Address = address
	v.Token = token

	return nil
}

// GetiFrame returns an HTML/CSS code to be used as an iFrame
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
		html, err = getTasksiFrame(v.Address, tasks, theme, apiURL, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Errorf("Couldn't create HTML code: %s", err.Error()))
			return
		}
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func getTasksiFrame(vikunjaAddress string, tasks []*Task, theme, apiURL string, limit int) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="TASKS-CONTAINER-BACKGROUND-COLOR">
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

        .tasks-container {
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

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = 'API-URL/v1/hash/vikunja?limit=API-LIMIT';
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

        fetchAndUpdate();
        
    </script>

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
                    <span class="info-label" style="color: {{ getTimeColor .DueDate }};"><i class="fa-solid fa-calendar-days"></i> Due: {{ .DueDate.Format "Jan 2, 2006" }}</span>
                {{ else if not .EndDate.IsZero }}
                    <span class="info-label" style="color: {{ getTimeColor .EndDate }};"><i class="fa-solid fa-calendar-days"></i> End: {{ .EndDate.Format "Jan 2, 2006" }}</span>
                {{ else if or (ne .RepeatAfter 0) (ne .RepeatMode 0) }}
                    {{ if or (eq .RepeatMode 0) (eq .RepeatMode 2) }}
                        <span class="info-label"><i class="fa-solid fa-calendar-days"></i> Repeats every {{ getRepeatAfter .RepeatAfter }}</span>
                    {{ else if eq .RepeatMode 1 }}
                        <span class="info-label"><i class="fa-solid fa-calendar-days"></i> Repeats monthly</span>
                    {{ end }}
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
	// Homarr theme
	scrollbarThumbBackgroundColor := "rgba(209, 219, 227, 1)"
	scrollbarTrackBackgroundColor := "#ffffff"
	if theme == "dark" {
		scrollbarThumbBackgroundColor = "#484d64"
		scrollbarTrackBackgroundColor = "rgba(37, 40, 53, 1)"
	}

	if apiURL != "" {
		html = strings.Replace(html, "API-URL", apiURL, -1)
		html = strings.Replace(html, "API-LIMIT", strconv.Itoa(limit), -1)
	} else {
		html = strings.Replace(html, "fetchAndUpdate();", "// fetchAndUpdate", -1)
	}

	html = strings.Replace(html, "VIKUNJA-ADDRESS", vikunjaAddress, -1)
	html = strings.Replace(html, "TASKS-CONTAINER-BACKGROUND-COLOR", theme, -1)
	html = strings.Replace(html, "TASKS-CONTAINER-BACKGROUND-IMAGE", backgroundImageURL, -1)
	html = strings.Replace(html, "SCROLLBAR-THUMB-BACKGROUND-COLOR", scrollbarThumbBackgroundColor, -1)
	html = strings.Replace(html, "SCROLLBAR-TRACK-BACKGROUND-COLOR", scrollbarTrackBackgroundColor, -1)

	divideFunc := template.FuncMap{
		"getRepeatAfter": func(a int) string {
			hours := float64(a) / 3600
			if hours != float64(int(hours)) {
				return fmt.Sprintf("%.1f hours", hours)
			}

			if hours < 24 {
				return fmt.Sprintf("%d hours", int(hours))
			}

			days := hours / 24
			if days != float64(int(days)) {
				return fmt.Sprintf("%d hours", int(hours))
			}

			return fmt.Sprintf("%d days", int(days))
		},
		"getTimeColor": func(t time.Time) string {
			if t.Before(time.Now()) {
				return "#ff4136"
			}
			if sources.IsToday(t) {
				return "#ff851b"
			}

			return "#99b6bb"
		},
	}

	tmpl := template.Must(template.New("tasks").Funcs(divideFunc).Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, tasks)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

// GetHash returns the hash of the tasks
func (v *Vikunja) GetHash(c *gin.Context) {
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

	pTasks := []*Task{}
	if limit != 0 {
		pTasks, err = v.GetTasks(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	var tasks []Task
	for _, task := range pTasks {
		tasks = append(tasks, *task)
	}

	hash := sources.GetHash(tasks)

	c.JSON(http.StatusOK, gin.H{"hash": fmt.Sprintf("%x", hash)})
}
