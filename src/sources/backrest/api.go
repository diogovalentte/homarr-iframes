package backrest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (b *Backrest) baseRequest(method, url string, body io.Reader, target interface{}) error {
	client := &http.Client{}
	if body == nil {
		body = strings.NewReader("{}")
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if b.username != "" && b.password != "" {
		req.SetBasicAuth(b.username, b.password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if string(resBody) == "Unauthorized (No Authorization Header)\n" {
			return fmt.Errorf("error: %s. Authentication is enabled. Set BACKREST_USERNAME and BACKREST_PASSWORD variables", resp.Status)
		}
		return fmt.Errorf("error: %s", resp.Status)
	}

	if err := json.Unmarshal(resBody, target); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\n. Reponse body: %s", err.Error(), string(resBody))
	}

	return nil
}
