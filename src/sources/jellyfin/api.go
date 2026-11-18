package jellyfin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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
				j.Address, item.ID, primaryTag)
		}

		if len(item.BackdropImageTags) > 0 {
			item.BackdropImageURL = fmt.Sprintf("%s/Items/%s/Images/Backdrop?tag=%s",
				j.Address, item.ID, item.BackdropImageTags[0])
		}

		item.ItemURL = fmt.Sprintf("%s/web/#/details?id=%s&serverId=%s", j.Address, item.ID, item.ServerId)
	}

	if limit > 0 && len(latestItems) > limit {
		latestItems = latestItems[:limit]
	}

	return latestItems, nil
}

// GetSessions returns the sessions that are currently playing something.
// Idle sessions and sessions of users in excludedUsers (lowercase usernames) are left out.
func (j *Jellyfin) GetSessions(limit, activeWithinSeconds int, excludedUsers map[string]bool) ([]*Session, error) {
	var jellyfinURL string = j.InternalAddress + "/Sessions"

	if activeWithinSeconds < 1 {
		activeWithinSeconds = 60
	}

	querySeconds := strconv.Itoa(activeWithinSeconds)
	queryParams := "?activeWithinSeconds=" + querySeconds

	jellyfinURL = jellyfinURL + queryParams

	var sessions []*Session
	if err := j.baseRequest(http.MethodGet, jellyfinURL, nil, &sessions); err != nil {
		return nil, fmt.Errorf("error getting sessions: %w", err)
	}

	playingSessions := make([]*Session, 0, len(sessions))
	for _, session := range sessions {
		if session.NowPlayingItem == nil {
			continue
		}
		if excludedUsers[strings.ToLower(session.UserName)] {
			continue
		}

		session.UserAvatarURL = fmt.Sprintf("%s/Users/%s/Images/Primary",
			j.Address, session.UserID)

		itemID := session.NowPlayingItem.ID

		if session.NowPlayingItem.SeriesID != "" {
			itemID = session.NowPlayingItem.SeriesID

			session.NowPlayingItem.BackdropImageURL = fmt.Sprintf("%s/Items/%s/Images/Backdrop",
				j.Address, itemID)

			session.NowPlayingItem.EpisodeURL = fmt.Sprintf("%s/web/#/details?id=%s&serverId=%s",
				j.Address, session.NowPlayingItem.ID, session.NowPlayingItem.ServerId)
		} else if len(session.NowPlayingItem.BackdropImageTags) > 0 {
			session.NowPlayingItem.BackdropImageURL = fmt.Sprintf("%s/Items/%s/Images/Backdrop?tag=%s",
				j.Address, itemID, session.NowPlayingItem.BackdropImageTags[0])
		}

		if primaryTag, hasPrimary := session.NowPlayingItem.ImageTags["Primary"]; hasPrimary {
			session.NowPlayingItem.PrimaryImageURL = fmt.Sprintf("%s/Items/%s/Images/Primary?tag=%s",
				j.Address, itemID, primaryTag)
		}

		session.NowPlayingItem.ItemURL = fmt.Sprintf("%s/web/#/details?id=%s&serverId=%s",
			j.Address, itemID, session.NowPlayingItem.ServerId)

		playingSessions = append(playingSessions, session)
	}
	sessions = playingSessions

	if limit > 0 && len(sessions) > limit {
		sessions = sessions[:limit]
	}

	return sessions, nil
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
