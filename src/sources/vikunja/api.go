package vikunja

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// GetTasks get not done tasks with using a custom ordering.
// Can also limit the number of tasks returned.
// The project_id is the id of the project to get the tasks from. Empty gets all tasks from all projects.
func (v *Vikunja) GetTasks(limit int, projectID int) ([]*Task, error) {
	target := []*Task{}

	path := "/api/v1/tasks/all?sort_by=due_date&order_by=asc&sort_by=end_date&order_by=asc&sort_by=created&order_by=desc&filter_by=done&filter_value=false&filter_comparator=equals"
	if limit > 0 {
		path = path + fmt.Sprintf("&per_page=%d", limit)
	}

	if projectID > 0 {
		path = path + fmt.Sprintf("&filter_concat=and&filter_by=project_id&filter_value=%d&filter_comparator=equals", projectID)
	}

	err := v.baseRequest("GET", v.Address+path, nil, &target)
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

func (v *Vikunja) SetTaskDone(taskId int) error {
	path := "/api/v1/tasks/" + strconv.Itoa(taskId)
	body := []byte(`{"done": true}`)
	task := &Task{}

	err := v.baseRequest("POST", v.Address+path, bytes.NewBuffer(body), task)
	if err != nil {
		return err
	}

	if task.Done != true {
		return fmt.Errorf("task not done")
	}

	return nil
}

func (v *Vikunja) GetProjects() (map[int]*Project, error) {
	path := "/api/v1/projects"
	projects := []*Project{}

	err := v.baseRequest("GET", v.Address+path, nil, &projects)
	if err != nil {
		return nil, err
	}

	projectsMap := make(map[int]*Project)
	for _, project := range projects {
		projectsMap[project.ID] = project
	}

	return projectsMap, nil
}

func (v *Vikunja) GetProject(projectId int) (*Project, error) {
	path := "/api/v1/projects/" + strconv.Itoa(projectId)
	project := &Project{}

	err := v.baseRequest("GET", v.Address+path, nil, &project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// SetInMemoryInstanceProjects sets the instanceProjects variable to the projects.
func (v *Vikunja) SetInMemoryInstanceProjects() error {
	projects, err := v.GetProjects()
	if err != nil {
		return err
	}

	instanceProjects = projects

	return nil
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
