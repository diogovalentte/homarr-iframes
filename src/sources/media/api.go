package media

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func baseRequest(method, url string, body io.Reader, target interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

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
		return fmt.Errorf("error unmarshaling JSON: %s\n. Reponse body: %s", err.Error(), string(resBody))
	}

	return nil
}

func getReleaseCoverImageURL(images []defaultReleaseImagesResponse) string {
	if len(images) == 0 {
		return ""
	}

	for _, image := range images {
		if image.CoverType == "poster" {
			return image.RemoteURL
		}
	}

	return images[0].RemoteURL
}

// isReleaseDateWithinDateRange checks if it's within a given date range.
// startDate is inclusive, endDate is exclusive.
func isReleaseDateWithinDateRange(releaseDate, startDate, endDate time.Time) bool {
	return (releaseDate.After(startDate) || releaseDate.Equal(startDate)) && releaseDate.Before(endDate)
}
