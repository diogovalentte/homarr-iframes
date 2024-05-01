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

	"github.com/diogovalentte/homarr-iframes/src/sources"
)

var (
	backgroundImageURL = "https://vikunja.io/images/vikunja.png"
	instanceProjects   = make(map[int]*Project)
)

var v *Vikunja

// Vikunja is the a source
type Vikunja struct {
	Address string
	Token   string
}

func New(address, token string) (*Vikunja, error) {
	if v != nil {
		return v, nil
	}

	newV := &Vikunja{}
	err := newV.Init(address, token)
	if err != nil {
		return nil, err
	}

	v = newV

	return v, nil
}

// Init sets the Vikunja properties from the configs
func (v *Vikunja) Init(address, token string) error {
	if address == "" || token == "" {
		return fmt.Errorf("VIKUNJA_ADDRESS and VIKUNJA_TOKEN variables should be set")
	}

	if strings.HasSuffix(address, "/") {
		address = address[:len(address)-1]
	}

	v.Address = address
	v.Token = token

	err := v.SetInMemoryInstanceProjects()
	if err != nil {
		return fmt.Errorf("couldn't get Vikunja projects: %s", err.Error())
	}

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

	var (
		showCreated  bool
		showDue      bool
		showPriority bool
		showProject  bool
	)
	showCreatedStr := c.Query("showCreated")
	if showCreatedStr == "" {
		showCreated = true
	} else {
		showCreated, err = strconv.ParseBool(showCreatedStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "showCreated must be a boolean"})
			return
		}
	}
	showDueStr := c.Query("showDue")
	if showDueStr == "" {
		showDue = true
	} else {
		showDue, err = strconv.ParseBool(showDueStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "showDue must be a boolean"})
			return
		}
	}
	showPriorityStr := c.Query("showPriority")
	if showPriorityStr == "" {
		showPriority = true
	} else {
		showPriority, err = strconv.ParseBool(showPriorityStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "showPriority must be a boolean"})
			return
		}
	}
	showProjectStr := c.Query("showProject")
	if showProjectStr == "" {
		showProject = true
	} else {
		showProject, err = strconv.ParseBool(showProjectStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "showProject must be a boolean"})
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
		html, err = getTasksiFrame(v.Address, tasks, theme, apiURL, limit, showCreated, showDue, showPriority, showProject)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Errorf("Couldn't create HTML code: %s", err.Error()))
			return
		}
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func getTasksiFrame(vikunjaAddress string, tasks []*Task, theme, apiURL string, limit int, showCreated, showDue, showPriority, showProject bool) ([]byte, error) {
	html := `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="referrer" content="no-referrer"> <!-- If not set, can't load some images when behind a domain or reverse proxy -->
    <meta name="color-scheme" content="{{ .Theme }}">
    <script src="https://kit.fontawesome.com/3f763b063a.js" crossorigin="anonymous"></script>
    <title>Vikunja iFrame</title>
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
            background-image: url('{{ .BackgroundImageURL }}');
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

        .set-task-done-container {
            display: inline-block;
            background-color: transparent;
            margin: 20px 20px 20px 10px;
            border-radius: 5px;
            width: 70px;
            text-align: center;
        }

        .set-task-done-button {
            color: white;
            background-color: #04c9b7;
            padding: 0.25rem 0.75rem;
            border-radius: 0.5rem;
            border: 1px solid rgb(4, 201, 183);
            font-weight: bold;
        }

        button.set-task-done-button:hover {
            filter: brightness(0.9)
        }
    </style>

    <script>
        let lastHash = null;

        async function fetchData() {
            try {
                var url = '{{ .APIURL }}/v1/hash/vikunja?limit={{ .APILimit }}';
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

    <script>
      function setTaskDone(taskId) {
        try {
            var xhr = new XMLHttpRequest();
            var url = '{{ .APIURL }}/v1/iframe/vikunja/set_task_done?taskId=' + encodeURIComponent(taskId);
            xhr.open('PATCH', url, true);
            xhr.setRequestHeader('Content-Type', 'application/json');

            xhr.onload = function () {
              if (xhr.status >= 200 && xhr.status < 300) {
                console.log('Request to set task ', taskId, ' as done finished with success:', xhr.responseText);
                location.reload();
              } else {
                console.log('Request to set task ', taskId, ' as done failed:', xhr.responseText);
                handleSetTaskDoneError("task-" + taskId)
              }
            };

            xhr.onerror = function () {
              console.log('Request to set task ', taskId, ' as done failed:', xhr.responseText);
              handleSetTaskDoneError("task-" + taskId)
            };

            xhr.send(null);
        } catch (error) {
            console.log('Request to set task ', taskId, ' as done failed:', xhr.responseText);
            handleSetTaskDoneError("task-" + taskId)
        }
      }

      function handleSetTaskDoneError(buttonId) {
        var button = document.getElementById(buttonId);
        button.textContent = "ERROR";
        button.style.backgroundColor = "red";
        button.style.borderColor = "red";
      }
    </script>

</head>
<body>
{{ range .Tasks }}
    <div class="tasks-container">

        <div class="background-image"></div>

        <div style="padding-left: 20px;" class="text-wrap">

            <a href="{{ with . }}{{ $.VikunjaAddress }}{{ end }}/tasks/{{ .ID }}" target="_blank" class="task-title">
                {{ with . }}{{ if $.ShowPriority }}
                    {{ if eq .Priority 3 }}
                        <span style="color: #ff851b;" class="info-label">! High</span>{{ .Title }}
                    {{ else if eq .Priority 4 }}
                        <span style="color: #ff4136;" class="info-label">! Urgent</span>{{ .Title }}
                    {{ else if eq .Priority 5 }}
                        <span style="color: #ff4136;" class="info-label">! DO NOW !</span>{{ .Title }}
                    {{ else }}
                        {{ .Title }}
                    {{ end }}
                {{ end }}{{ end }}
            </a>

            <div>

                {{ with . }}{{ if $.ShowCreated }}
                    <span class="info-label"><i class="fa-solid fa-calendar-days"></i> {{ .CreatedAt.Format "Jan 2, 2006" }}</span>
                {{ end }}{{ end }}
            
                {{ with . }}{{ if $.ShowDue }}
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
                {{ end }}{{ end }}

                {{ with . }}{{ if $.ShowProject }}
                    {{ if gt .ProjectID 1 }} <!-- 1 = Inbox; -1 = Favorites   -->
                        {{ $project := getTaskProject .ProjectID }}
                        {{ if $project.Title }}
                            <span class="info-label" style="color: #{{ $project.HexColor }};"><i class="fa-solid fa-layer-group"></i> {{ $project.Title }}</span>
                        {{ end }}
                    {{ end }}
                {{ end }}{{ end }}

            </div>

        </div>

        {{ with . }}{{ if $.APIURL }}
            <div class="set-task-done-container">
                <button id="task-{{ .ID }}" onclick="setTaskDone('{{ .ID }}')" class="set-task-done-button" onmouseenter="this.style.cursor='pointer';">Done</button>
            </div>
        {{ end }}{{ end }}

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
		Tasks:                         tasks,
		Theme:                         theme,
		APIURL:                        apiURL,
		APILimit:                      limit,
		VikunjaAddress:                vikunjaAddress,
		BackgroundImageURL:            backgroundImageURL,
		ScrollbarThumbBackgroundColor: scrollbarThumbBackgroundColor,
		ScrollbarTrackBackgroundColor: scrollbarTrackBackgroundColor,
		ShowCreated:                   showCreated,
		ShowDue:                       showDue,
		ShowPriority:                  showPriority,
		ShowProject:                   showProject,
	}

	templateFuncs := template.FuncMap{
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
		"getTaskProject": func(projectID int) *Project {
			project, ok := instanceProjects[projectID]
			if ok {
				return project
			}

			v, err := New("", "")
			if err != nil {
				return &Project{}
			}
			project, err = v.GetProject(projectID)
			if err != nil {
				return &Project{}
			}
			return project
		},
	}

	tmpl := template.Must(template.New("tasks").Funcs(templateFuncs).Parse(html))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, &templateData)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

type iframeTemplateData struct {
	Tasks                         []*Task
	Theme                         string
	APIURL                        string
	APILimit                      int
	VikunjaAddress                string
	BackgroundImageURL            string
	ScrollbarThumbBackgroundColor string
	ScrollbarTrackBackgroundColor string
	ShowCreated                   bool
	ShowDue                       bool
	ShowPriority                  bool
	ShowProject                   bool
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
