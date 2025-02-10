package vikunja

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// GetTasks get not done tasks with using a custom ordering.
// Can also limit the number of tasks returned.
// The project_id is the id of the project to get the tasks from. Empty gets all tasks from all projects.
func (v *Vikunja) GetTasks(limit int, projectID int, excludeProjectIDs []*int) ([]*Task, error) {
	var tasks []*Task
	var err error
	if isGreater, err := v.IsVikunjaVersionGreaterOrEqualTo("0.24.0"); err != nil {
		return nil, err
	} else if isGreater {
		tasks, err = v.getTasksV2(limit, projectID)
	} else {
		tasks, err = v.getTasksV1(limit, projectID)
	}

	if err != nil {
		return nil, err
	}

	if len(excludeProjectIDs) > 0 {
		var filteredTasks []*Task
		for _, task := range tasks {
			exclude := false
			for _, excludeProjectID := range excludeProjectIDs {
				if *excludeProjectID == -1 {
					if task.IsFavorite {
						exclude = true
						break
					}
				} else {
					if task.ProjectID == *excludeProjectID {
						exclude = true
						break
					}
				}
			}
			if !exclude {
				filteredTasks = append(filteredTasks, task)
			}
		}

		return filteredTasks, nil
	}

	return tasks, nil
}

func (v *Vikunja) getTasksV2(limit int, projectID int) ([]*Task, error) {
	target := []*Task{}

	var path string
	if projectID > 0 {
		path = fmt.Sprintf("/api/v1/projects/%d/tasks", projectID)
	} else {
		path = "/api/v1/tasks/all"
	}
	path = path + "?sort_by=due_date&order_by=asc&sort_by=end_date&order_by=asc&sort_by=priority&order_by=desc&sort_by=created&order_by=desc&filter=done=false"
	if limit > 0 {
		path = path + fmt.Sprintf("&per_page=%d", limit)
	}

	err := v.baseRequest("GET", v.InternalAddress+path, nil, &target)
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	if projectID == -1 {
		for _, task := range target {
			if task.IsFavorite {
				tasks = append(tasks, task)
			}
		}
	} else {
		tasks = target
	}

	return tasks, nil
}

func (v *Vikunja) getTasksV1(limit int, projectID int) ([]*Task, error) {
	target := []*Task{}

	path := "/api/v1/tasks/all?sort_by=due_date&order_by=asc&sort_by=end_date&order_by=asc&sort_by=priority&order_by=desc&sort_by=created&order_by=desc&filter_by=done&filter_value=false&filter_comparator=equals"
	if limit > 0 {
		path = path + fmt.Sprintf("&per_page=%d", limit)
	}

	if projectID > 0 {
		path = path + fmt.Sprintf("&filter_concat=and&filter_by=project_id&filter_value=%d&filter_comparator=equals", projectID)
	}

	err := v.baseRequest("GET", v.InternalAddress+path, nil, &target)
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	if projectID == -1 {
		for _, task := range target {
			if task.IsFavorite {
				tasks = append(tasks, task)
			}
		}
	} else {
		tasks = target
	}

	return tasks, nil
}

func (v *Vikunja) SetTaskDone(taskID int) error {
	path := "/api/v1/tasks/" + strconv.Itoa(taskID)
	body := []byte(`{"done": true}`)
	task := &Task{}

	err := v.baseRequest("POST", v.InternalAddress+path, bytes.NewBuffer(body), task)
	if err != nil {
		return err
	}

	if !task.Done {
		return fmt.Errorf("task not done")
	}

	return nil
}

func (v *Vikunja) GetProjects() ([]*Project, error) {
	path := "/api/v1/projects"
	projects := []*Project{}

	err := v.baseRequest("GET", v.InternalAddress+path, nil, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (v *Vikunja) GetProject(projectID int) (*Project, error) {
	path := "/api/v1/projects/" + strconv.Itoa(projectID)
	project := &Project{}

	err := v.baseRequest("GET", v.InternalAddress+path, nil, &project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (v *Vikunja) baseRequest(method, url string, body io.Reader, target interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+v.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", resp.Status)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(resBody, target); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\nreponse text: %s", err.Error(), string(resBody))
	}

	return nil
}

// IsVersionGreaterOrEqualToVikunjaVersion checks if version is greater or equal to the Vikunja version.
// version1 and version2 should be like "1.25.3"
func (v *Vikunja) IsVikunjaVersionGreaterOrEqualTo(version string) (bool, error) {
	vikunjaVersion, err := v.getVikunjaVersion()
	if err != nil {
		return false, err
	}

	return IsVersionGreaterOrEqualTo(vikunjaVersion, version)
}

// IsVersionGreaterOrEqualTo checks if version1 is greater or equal to version2.
// version1 and version2 should be like "1.25.3"
func IsVersionGreaterOrEqualTo(version1 string, version2 string) (bool, error) {
	version1Parts := strings.Split(version1, ".")
	version2Parts := strings.Split(version2, ".")

	if len(version1Parts) != 3 || len(version2Parts) != 3 {
		return false, fmt.Errorf("version is not in the correct format")
	}

	for i := 0; i < 3; i++ {
		version1Part, err := strconv.Atoi(version1Parts[i])
		if err != nil {
			return false, err
		}
		version2Part, err := strconv.Atoi(version2Parts[i])
		if err != nil {
			return false, err
		}

		if version2Part > version1Part {
			return false, nil
		} else if version1Part > version2Part {
			return true, nil
		}
	}

	return true, nil
}

func (v *Vikunja) getVikunjaVersion() (string, error) {
	info := &Info{}
	err := v.baseRequest("GET", v.InternalAddress+"/api/v1/info", nil, info)
	if err != nil {
		return "", err
	}

	version := strings.Replace(info.Version, "v", "", 1)
	if version == "" {
		return "", fmt.Errorf("version is empty")
	}

	return version, nil
}
