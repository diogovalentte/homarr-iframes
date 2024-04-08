package vikunja

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetToken returns logs in the API and return a token
func (v *Vikunja) GetToken() (string, error) {
	data := loginRequestBody{
		Username: v.Username,
		Password: v.Password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("Login Error: error while creating request body: %w", err)
	}

	resp, err := http.Post(v.Address+"/api/v1/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Login Error: error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Login Error: error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Login Error: unexpected status code: %s\nResponse text: %s", resp.Status, string(body))
	}

	target := loginResponseBody{}
	if err := json.Unmarshal(body, &target); err != nil {
		return "", fmt.Errorf("Login Error: error unmarshaling JSON: %s\nReponse text: %s", err.Error(), string(body))
	}

	return target.Token, nil
}

// GetTasks get not done tasks with using a custom ordering.
// Can also limit the number of tasks returned.
func (v *Vikunja) GetTasks(limit int) ([]*Task, error) {
	target := []*Task{}

	path := "/api/v1/tasks/all?sort_by=due_date&order_by=asc&sort_by=end_date&order_by=asc&sort_by=created&order_by=desc&filter_by=done&filter_value=false&filter_comparator=equals"
	if limit > 0 {
		path = path + fmt.Sprintf("&per_page=%d", limit)
	}
	err := v.baseRequest(v.Address+path, &target)
	if err != nil {
		return nil, err
	}

	return target, err
}

func (v *Vikunja) baseRequest(url string, target interface{}) error {
	token, err := v.GetToken()
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\nreponse text: %s", err.Error(), string(body))
	}

	return nil
}
