package openarchiver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (o *OpenArchiver) GetIngestionSources(limit int) ([]IngestionSource, error) {
	var sources []IngestionSource
	err := o.baseRequest(http.MethodGet, o.InternalAddress+"/api/v1/ingestion-sources", nil, &sources)
	if err != nil {
		return nil, fmt.Errorf("error getting ingestion sources: %w", err)
	}

	if limit > 0 && len(sources) > limit {
		sources = sources[:limit]
	}

	return sources, nil
}

func (o *OpenArchiver) baseRequest(method, url string, body io.Reader, target any) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+o.SuperAPIKey)
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

	if err := json.Unmarshal(resBody, target); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\nreponse text: %s", err.Error(), string(resBody))
	}

	return nil
}
