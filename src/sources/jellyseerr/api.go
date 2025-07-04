package jellyseerr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/sources/overseerr"
)

func (j *Jellyseerr) GetRequests(limit int, filter, sort string, requestedBy int) ([]overseerr.Request, error) {
	if limit == 0 {
		return []overseerr.Request{}, nil
	}
	path := fmt.Sprintf("/api/v1/request?take=%d", limit)
	if filter != "" {
		path += fmt.Sprintf("&filter=%s", filter)
	}
	if sort != "" {
		path += fmt.Sprintf("&sort=%s", sort)
	}
	if requestedBy > 0 {
		path += fmt.Sprintf("&requestedBy=%d", requestedBy)
	}

	var responseData overseerr.GetRequestsResponse
	if err := j.baseRequest(http.MethodGet, j.InternalAddress+path, nil, &responseData); err != nil {
		return nil, fmt.Errorf("error getting requests: %w", err)
	}

	return responseData.Requests, nil
}

func (j *Jellyseerr) GetMedia(limit int, filter, sort string) ([]overseerr.Media, error) {
	if limit == 0 {
		return []overseerr.Media{}, nil
	}
	path := fmt.Sprintf("/api/v1/media?take=%d", limit)
	if filter != "" {
		path += fmt.Sprintf("&filter=%s", filter)
	}
	if sort != "" {
		path += fmt.Sprintf("&sort=%s", sort)
	}

	var responseData overseerr.GetMediaResponse
	if err := j.baseRequest(http.MethodGet, j.InternalAddress+path, nil, &responseData); err != nil {
		return nil, fmt.Errorf("error getting media: %w", err)
	}

	return responseData.Media, nil
}

func (j *Jellyseerr) GetMovie(id int) (overseerr.GenericMedia, error) {
	var responseData getMovieResponse
	if err := j.baseRequest(http.MethodGet, j.InternalAddress+"/api/v1/movie/"+fmt.Sprint(id), nil, &responseData); err != nil {
		return overseerr.GenericMedia{}, fmt.Errorf("error getting movie: %w", err)
	}

	movie := overseerr.GenericMedia{
		Name:         responseData.Title,
		ID:           responseData.ID,
		ReleaseDate:  responseData.ReleaseDate,
		PosterPath:   responseData.PosterPath,
		BackdropPath: responseData.BackdropPath,
	}

	return movie, nil
}

type getMovieResponse struct {
	Title        string `json:"originalTitle"`
	ReleaseDate  string `json:"releaseDate"`
	BackdropPath string `json:"backdropPath"`
	PosterPath   string `json:"posterPath"`
	ID           int    `json:"id"`
}

func (j *Jellyseerr) GetTV(id int) (overseerr.GenericMedia, error) {
	var responseData getTvResponse
	if err := j.baseRequest(http.MethodGet, j.InternalAddress+"/api/v1/tv/"+fmt.Sprint(id), nil, &responseData); err != nil {
		return overseerr.GenericMedia{}, fmt.Errorf("error getting tv show: %w", err)
	}

	tvShow := overseerr.GenericMedia{
		Name:         responseData.Name,
		ID:           responseData.ID,
		ReleaseDate:  responseData.FirstAirDate,
		PosterPath:   responseData.PosterPath,
		BackdropPath: responseData.BackdropPath,
	}

	return tvShow, nil
}

type getTvResponse struct {
	Name         string `json:"originalName"`
	FirstAirDate string `json:"firstAirDate"`
	BackdropPath string `json:"backdropPath"`
	PosterPath   string `json:"posterPath"`
	ID           int    `json:"id"`
}

func (j *Jellyseerr) baseRequest(method, url string, body io.Reader, target any) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-Api-Key", j.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request status: %s", resp.Status)
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

func (j *Jellyseerr) GetIframeData(limit int, filter, sort string, requestedBy int, getMedia bool) ([]overseerr.IframeRequestData, error) {
	iframeData := []overseerr.IframeRequestData{}

	if !getMedia {
		requests, err := j.GetRequests(limit, filter, sort, requestedBy)
		if err != nil {
			return nil, err
		}

		for _, request := range requests {
			var data overseerr.IframeRequestData
			err = j.setMediaData(&request.Media, &data)
			if err != nil {
				return nil, fmt.Errorf("error setting media data for request %d: %w", request.ID, err)
			}
			data.Request.Username = request.RequestedBy.Username
			if strings.HasPrefix(request.RequestedBy.Avatar, "/avatarproxy/") {
				data.Request.AvatarURL = j.Address + request.RequestedBy.Avatar
			} else {
				data.Request.AvatarURL = request.RequestedBy.Avatar
			}
			data.Request.UserProfileURL = fmt.Sprintf("%s/users/%d", j.Address, request.RequestedBy.ID)
			data.Request.UserID = request.RequestedBy.ID
			data.Status = getRequestStatusName(request.Status, request.Media.Status)

			iframeData = append(iframeData, data)
		}
	} else {
		media, err := j.GetMedia(limit, filter, sort)
		if err != nil {
			return nil, err
		}
		for _, m := range media {
			var data overseerr.IframeRequestData
			err = j.setMediaData(&m, &data)
			if err != nil {
				return nil, fmt.Errorf("error setting media data: %w", err)
			}
			data.Status = getRequestStatusName(2, m.Status) // Default status for media

			iframeData = append(iframeData, data)
		}
	}

	return iframeData, nil
}

func (j *Jellyseerr) setMediaData(media *overseerr.Media, iframe *overseerr.IframeRequestData) error {
	if media.TMDBID == 0 {
		return fmt.Errorf("invalid TMDB ID for request %d", media.ID)
	}

	var err error
	var mediaInfo overseerr.GenericMedia
	switch media.Type {
	case "movie":
		mediaInfo, err = j.GetMovie(media.TMDBID)
	case "tv":
		mediaInfo, err = j.GetTV(media.TMDBID)
	default:
		return fmt.Errorf("invalid media type: %s", media.Type)
	}
	if err != nil {
		return err
	}

	iframe.Media.Name = mediaInfo.Name
	iframe.Media.Type = media.Type
	iframe.Media.TMDBID = media.TMDBID
	iframe.Media.Year = strings.Split(mediaInfo.ReleaseDate, "-")[0]
	iframe.Media.URL = fmt.Sprintf("%s/%s/%d", j.Address, media.Type, media.TMDBID)
	if mediaInfo.PosterPath != "" {
		iframe.Media.PosterURL = overseerr.TMDBPosterImageBasePath + mediaInfo.PosterPath
	} else {
		iframe.Media.PosterURL = overseerr.DefaultBackgroundImageURL
	}

	if mediaInfo.BackdropPath != "" {
		iframe.Media.BackdropURL = overseerr.TMDBBackdropImageBasePath + mediaInfo.BackdropPath
	} else {
		iframe.Media.BackdropURL = iframe.Media.PosterURL
	}

	return nil
}

// getRequestStatusName returns the HTML/CSS properties of the request status
// to be used in the iframe.
func getRequestStatusName(reqStatus, mediaStatus int) overseerr.IframeStatus {
	var status overseerr.IframeStatus
	switch reqStatus {
	case 1: // Pending
		status.Status = "Pending"
		status.Color = "#fe99ff"
		status.BackgroundColor = "#f000e733"
	case 2, 5: // Approved, Completed
		switch mediaStatus {
		case 1:
			status.Status = "Unkown"
			status.Color = "#99fff2"
			status.BackgroundColor = "#00f0dc33"
		case 2:
			status.Status = "Pending"
			status.Color = "#fe99ff"
			status.BackgroundColor = "#f000e733"
		case 3:
			status.Status = "Approved"
			status.Color = "#d0bfff"
			status.BackgroundColor = "#6741d933"
		case 4:
			status.Status = "Partial"
			status.Color = "#ff9f1a"
			status.BackgroundColor = "#f08c0033"
		case 5:
			status.Status = "Available"
			status.Color = "#b2f2bb"
			status.BackgroundColor = "#2f9e4433"
		case 7:
			status.Status = "Deleted"
			status.Color = "#f2b2ba"
			status.BackgroundColor = "#9e302f33"
		default:
			status.Status = "Approved"
			status.Color = "#d0bfff"
			status.BackgroundColor = "#6741d933"
		}
	case 3: // Declined
		status.Status = "Declined"
		status.Color = "#f2b2ba"
		status.BackgroundColor = "#9e302f33"
	case 4: // Failed
		status.Status = "Failed"
		status.Color = "#f2b2ba"
		status.BackgroundColor = "#9e302f33"
	default:
		status.Status = "Unkown"
		status.Color = "#99fff2"
		status.BackgroundColor = "#00f0dc33"
	}

	return status
}
