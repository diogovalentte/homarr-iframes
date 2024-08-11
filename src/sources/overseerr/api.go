package overseerr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (o *Overseerr) GetRequests(limit int, filter, sort string, requestedBy int) ([]Request, error) {
	if limit == 0 {
		return []Request{}, nil
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

	var responseData getRequestsResponse
	if err := o.baseRequest(http.MethodGet, o.InternalAddress+path, nil, &responseData); err != nil {
		return nil, fmt.Errorf("error getting requests: %w", err)
	}

	return responseData.Requests, nil
}

type getRequestsResponse struct {
	Requests []Request `json:"results"`
}

func (o *Overseerr) GetMedia(mediaType string, tmdbID int) (GenericMedia, error) {
	if tmdbID == 0 {
		return GenericMedia{}, fmt.Errorf("invalid TMDB ID")
	}
	var err error
	var media GenericMedia
	if mediaType == "movie" {
		media, err = o.GetMovie(tmdbID)
	} else if mediaType == "tv" {
		media, err = o.GetTv(tmdbID)
	} else {
		return GenericMedia{}, fmt.Errorf("invalid media type")
	}

	if err != nil {
		return GenericMedia{}, fmt.Errorf("error getting media: %w", err)
	}

	return media, nil
}

func (o *Overseerr) GetMovie(id int) (GenericMedia, error) {
	var responseData getMovieResponse
	if err := o.baseRequest(http.MethodGet, o.InternalAddress+"/api/v1/movie/"+fmt.Sprint(id), nil, &responseData); err != nil {
		return GenericMedia{}, fmt.Errorf("error getting movie: %w", err)
	}

	movie := GenericMedia{
		Name:        responseData.Title,
		ID:          responseData.ID,
		ReleaseDate: responseData.ReleaseDate,
		PosterPath:  responseData.PosterPath,
	}

	return movie, nil
}

type getMovieResponse struct {
	Title       string `json:"originalTitle"`
	ID          int    `json:"id"`
	ReleaseDate string `json:"releaseDate"`
	PosterPath  string `json:"posterPath"`
}

func (o *Overseerr) GetTv(id int) (GenericMedia, error) {
	var responseData getTvResponse
	if err := o.baseRequest(http.MethodGet, o.InternalAddress+"/api/v1/tv/"+fmt.Sprint(id), nil, &responseData); err != nil {
		return GenericMedia{}, fmt.Errorf("error getting tv show: %w", err)
	}

	tvShow := GenericMedia{
		Name:        responseData.Name,
		ID:          responseData.ID,
		ReleaseDate: responseData.FirstAirDate,
		PosterPath:  responseData.PosterPath,
	}

	return tvShow, nil
}

type getTvResponse struct {
	Name         string `json:"originalName"`
	ID           int    `json:"id"`
	FirstAirDate string `json:"firstAirDate"`
	PosterPath   string `json:"posterPath"`
}

func (o *Overseerr) baseRequest(method, url string, body io.Reader, target interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-Api-Key", o.Token)
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

func (o *Overseerr) GetIframeData(limit int, filter, sort string, requestedBy int) ([]iframeRequestData, error) {
	requests, err := o.GetRequests(limit, filter, sort, requestedBy)
	if err != nil {
		return nil, err
	}

	iframeData := []iframeRequestData{}
	for _, request := range requests {
		media, err := o.GetMedia(request.Media.Type, request.Media.TMDBID)
		if err != nil {
			return nil, err
		}
		var data iframeRequestData
		data.Media.Name = media.Name
		data.Media.Type = request.Media.Type
		data.Media.TMDBID = request.Media.TMDBID
		data.Media.Year = strings.Split(media.ReleaseDate, "-")[0]
		data.Media.PosterURL = tmdbImageBasePath + media.PosterPath
		data.Request.Username = request.RequestedBy.Username
		data.Request.AvatarURL = request.RequestedBy.Avatar
		data.Request.UserID = request.RequestedBy.ID
		data.Status = getRequestStatusName(request.Status, request.Media.Status)

		iframeData = append(iframeData, data)
	}

	return iframeData, nil
}

type iframeRequestData struct {
	Media struct {
		Name      string
		Type      string
		TMDBID    int
		Year      string
		PosterURL string
	}
	Request struct {
		Username  string
		AvatarURL string
		UserID    int
	}
	Status iframeStatus
}

type iframeStatus struct {
	Status          string
	Color           string
	BackgroundColor string
}

// getRequestStatusName returns the HTML/CSS properties of the request status
// to be used in the iframe.
func getRequestStatusName(reqStatus, mediaStatus int) iframeStatus {
	var status iframeStatus
	switch reqStatus {
	case 1:
		status.Status = "Pending"
		status.Color = "#fe99ff"
		status.BackgroundColor = "#f000e733"
	case 2:
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
		default:
			status.Status = "Approved"
			status.Color = "#d0bfff"
			status.BackgroundColor = "#6741d933"
		}
	case 3:
		status.Status = "Declined"
		status.Color = "#f2b2ba"
		status.BackgroundColor = "#9e302f33"
	default:
		status.Status = "Unkown"
		status.Color = "#99fff2"
		status.BackgroundColor = "#00f0dc33"
	}

	return status
}
