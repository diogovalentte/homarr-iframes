package jellyfin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (j *JellyfinRecently) GetLatestItems(limit int, userId string, parentId string, includeItemTypes string, queryLimit int) ([]*Item, error) {
	if userId == "" {
		if j.UserId == "" {
			return nil, fmt.Errorf("userId is required and not provided")
		}
		userId = j.UserId
	}

	var jellyfinURL string = j.InternalAddress + "/Users/" + userId + "/Items/Latest"

	queryParams := "?api_key=" + j.APIKey

	if queryLimit > 0 {
		queryParams += "&Limit=" + strconv.Itoa(queryLimit)
	} else {
		queryParams += "&Limit=100"
	}

	if parentId != "" {
		queryParams += "&ParentId=" + parentId
	}

	if includeItemTypes != "" {
		queryParams += "&IncludeItemTypes=" + includeItemTypes
	} else {
		queryParams += "&IncludeItemTypes=Movie,Episode"
	}

	jellyfinURL = jellyfinURL + queryParams

	resp, err := http.Get(jellyfinURL)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", resp.Status)
	}

	var latestItems []*Item
	if err := json.NewDecoder(resp.Body).Decode(&latestItems); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
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
