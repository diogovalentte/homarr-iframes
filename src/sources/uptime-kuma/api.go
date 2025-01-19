package uptimekuma

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetStatusPageLastUpDownCount returns the number of up and down sites for the last heartbeat of a status page
func (u *UptimeKuma) GetStatusPageLastUpDownCount(slug string) (*UpDownSites, error) {
	target := HeartbeatResponse{}

	path := "/api/status-page/heartbeat/" + slug

	err := u.baseRequest(u.Address+path, &target)
	if err != nil {
		return &UpDownSites{}, err
	}

	upDownSites := &UpDownSites{}

	for _, site := range target.HeartbeatList {
		if len(site) == 0 {
			continue
		}
		lastHeartbeat := site[len(site)-1]
		if lastHeartbeat.Status == 1 {
			upDownSites.Up++
		} else {
			upDownSites.Down++
		}
	}

	return upDownSites, nil
}

func (u *UptimeKuma) baseRequest(url string, target interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

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
