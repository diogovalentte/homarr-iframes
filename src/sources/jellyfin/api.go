package jellyfin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (j *Jellyfin) GetLatestItems(limit, queryLimit int, userId, parentId, includeItemTypes string) ([]*Item, error) {
	if userId == "" {
		if j.userId == "" {
			return nil, fmt.Errorf("userId is required and not provided")
		}
		userId = j.userId
	}

	var jellyfinURL string = j.InternalAddress + "/Users/" + userId + "/Items/Latest"

	if queryLimit < 1 {
		queryLimit = 100
	}
	queryLimitStr := strconv.Itoa(queryLimit)
	queryParams := "?Limit=" + queryLimitStr

	if parentId != "" {
		queryParams += "&ParentId=" + parentId
	}

	if includeItemTypes == "" {
		includeItemTypes = "Movie,Episode"
	}
	queryParams += "&IncludeItemTypes=" + includeItemTypes

	jellyfinURL = jellyfinURL + queryParams

	var latestItems []*Item
	if err := j.baseRequest(http.MethodGet, jellyfinURL, nil, &latestItems); err != nil {
		return nil, fmt.Errorf("error getting latest items: %w", err)
	}

	for _, item := range latestItems {
		if primaryTag, hasPrimary := item.ImageTags["Primary"]; hasPrimary {
			item.PrimaryImageURL = fmt.Sprintf("%s/Items/%s/Images/Primary?tag=%s",
				j.InternalAddress, item.ID, primaryTag)
		}

		if len(item.BackdropImageTags) > 0 {
			item.BackdropImageURL = fmt.Sprintf("%s/Items/%s/Images/Backdrop?tag=%s",
				j.InternalAddress, item.ID, item.BackdropImageTags[0])
		}

		item.ItemURL = fmt.Sprintf("%s/web/#/details?id=%s&serverId=%s", j.Address, item.ID, item.ServerId)
	}

	if limit > 0 && len(latestItems) > limit {
		latestItems = latestItems[:limit]
	}

	return latestItems, nil
}

func (j *Jellyfin) baseRequest(method, url string, body io.Reader, target any) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("MediaBrowser Token=\"%s\"", j.APIKey))
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
