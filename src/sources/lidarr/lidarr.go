package lidarr

import (
	"fmt"
	"strings"
	"time"

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources/radarr"
)

var (
	l                  *Lidarr
	BackgroundImageURL = "https://avatars.githubusercontent.com/u/25025331"
)

type Lidarr struct {
	Address         string
	InternalAddress string
	APIKey          string
}

func New() (*Lidarr, error) {
	if l != nil {
		return l, nil
	}

	newL := &Lidarr{}
	err := newL.Init()
	if err != nil {
		return nil, err
	}

	l = newL

	return l, nil
}

func (l *Lidarr) Init() error {
	address, internalAddress, APIKey := config.GlobalConfigs.Lidarr.Address, config.GlobalConfigs.Lidarr.InternalAddress, config.GlobalConfigs.Lidarr.APIKey
	if address == "" || APIKey == "" {
		return fmt.Errorf("LIDARR_ADDRESS and LIDARR_API_KEY variables should be set")
	}

	l.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		l.InternalAddress = l.Address
	} else {
		l.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	l.APIKey = APIKey

	return nil
}

// GetCalendar returns the calendar of releases where the release date is between startDate and endDate.
// To get the calendar of a specific day, set the startDate and endDate to the same day.
// To get the calendar of a specific week, set the startDate to the first day of the week and endDate to the last day of the week.
// It considers only the date, not the time, so it'll get all releases that are released on that day.
func (l *Lidarr) GetCalendar(unmonitored bool, startDate, endDate time.Time) ([]*GetLidarrCalendarEntryResponse, error) {
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), endDate.Location())

	var entries []*GetLidarrCalendarEntryResponse
	err := baseRequest("GET", fmt.Sprintf("%s/api/v1/calendar?start=%s&end=%s&unmonitored=%v&includeArtist=true", l.InternalAddress, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), unmonitored), nil, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

type GetLidarrCalendarEntryResponse struct {
	Title          string `json:"title"`
	ForeignAlbumID string `json:"foreignAlbumId"`
	AlbumType      string `json:"albumType"`
	ReleaseDate    string `json:"releaseDate"`
	Artist         struct {
		ArtistName      string                                `json:"artistName"`
		ForeignArtistID string                                `json:"foreignArtistId"`
		Images          []radarr.DefaultReleaseImagesResponse `json:"images"`
	}
	Images     []radarr.DefaultReleaseImagesResponse `json:"images"`
	Statistics struct {
		TrackFileCount  int `json:"trackFileCount"`
		TotalTrackCount int `json:"totalTrackCount"`
	}
}

func (r *Lidarr) GetHealth() ([]*HealthEntry, error) {
	var entries []*HealthEntry
	err := baseRequest("GET", fmt.Sprintf("%s/api/v1/health", r.InternalAddress), nil, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

type HealthEntry struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	WikiURL string `json:"wikiUrl"`
}
