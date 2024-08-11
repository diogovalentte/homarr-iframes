package netdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

func (n *Netdata) GetAlarms(limit int) ([]Alarm, error) {
	if limit == 0 {
		return []Alarm{}, nil
	}

	path := "/api/v1/alarms"

	var resp *getAlarmsResponse
	err := n.baseRequest("GET", n.InternalAddress+path, nil, &resp)
	if err != nil {
		return nil, err
	}

	alarmsSliceSize := limit
	if len(resp.Alarms) < limit || limit < 0 {
		alarmsSliceSize = len(resp.Alarms)
	}
	alarms := make([]Alarm, 0, alarmsSliceSize)
	counter := 0
	for _, alarm := range resp.Alarms {
		alarm.LastStatusChange = time.Unix(alarm.LastStatusChangeUnix, 0)
		alarms = append(alarms, alarm)
		if limit > 0 {
			counter++
			if counter == limit {
				break
			}
		}
	}

	sort.Slice(alarms, func(i, j int) bool {
		return alarms[i].LastStatusChange.After(alarms[j].LastStatusChange)
	})

	return alarms, nil
}

type getAlarmsResponse struct {
	Alarms map[string]Alarm `json:"alarms"`
}

func (n *Netdata) baseRequest(method, url string, body io.Reader, target interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+n.Token)
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
