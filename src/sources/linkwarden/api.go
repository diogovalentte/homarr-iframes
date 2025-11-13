package linkwarden

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (l *Linkwarden) GetLinks(limit int, collectionID string) ([]*Link, error) {
	var linkwardenURL string
	if collectionID != "" {
		linkwardenURL = l.InternalAddress + "/api/v1/links?collectionId=" + collectionID
	} else {
		linkwardenURL = l.InternalAddress + "/api/v1/links"
	}

	linksResp := map[string][]*Link{}
	err := l.baseRequest(http.MethodGet, linkwardenURL, nil, &linksResp)
	if err != nil {
		return nil, fmt.Errorf("error while doing API request: %w", err)
	}

	links, exists := linksResp["response"]
	if !exists {
		return nil, fmt.Errorf("no 'response' field in API response")
	}

	if limit >= 0 {
		links = links[:limit]
	}

	return links, nil
}

func (l *Linkwarden) DeleteLink(linkId string) error {
	linkwardenURL := l.InternalAddress + "/api/v1/links/" + linkId

	err := l.baseRequest(http.MethodDelete, linkwardenURL, nil, nil)
	if err != nil {
		return fmt.Errorf("error while doing API request: %w", err)
	}

	return nil
}

func (l *Linkwarden) baseRequest(method, url string, body io.Reader, target any) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+l.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request status (%s): %s", resp.Status, string(resBody))
	}

	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if target != nil {
		if err := json.Unmarshal(resBody, target); err != nil {
			return fmt.Errorf("error unmarshaling JSON: %s\nreponse text: %s", err.Error(), string(resBody))
		}
	}

	return nil
}
