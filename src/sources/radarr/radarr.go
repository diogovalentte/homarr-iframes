package radarr

import (
	"fmt"
	"strings"
	"time"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var (
	r                  *Radarr
	BackgroundImageURL = "https://avatars.githubusercontent.com/u/25025331"
)

type Radarr struct {
	Address         string
	InternalAddress string
	APIKey          string
}

func New() (*Radarr, error) {
	if r != nil {
		return r, nil
	}

	newR := &Radarr{}
	err := newR.Init()
	if err != nil {
		return nil, err
	}

	r = newR

	return r, nil
}

func (r *Radarr) Init() error {
	address, internalAddress, APIKey := config.GlobalConfigs.Radarr.Address, config.GlobalConfigs.Radarr.InternalAddress, config.GlobalConfigs.Radarr.APIKey
	if address == "" || APIKey == "" {
		return fmt.Errorf("RADARR_ADDRESS and RADARR_API_KEY variables should be set")
	}

	r.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		r.InternalAddress = r.Address
	} else {
		r.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	r.APIKey = APIKey

	return nil
}

// GetCalendar returns the calendar of releases where the release date is between startDate and endDate.
// To get the calendar of a specific day, set the startDate and endDate to the same day.
// To get the calendar of a specific week, set the startDate to the first day of the week and endDate to the last day of the week.
// It considers only the date, not the time, so it'll get all releases that are released on that day.
func (r *Radarr) GetCalendar(unmonitored bool, startDate, endDate time.Time) ([]*getRadarrCalendarEntryResponse, error) {
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), endDate.Location())

	var entries []*getRadarrCalendarEntryResponse
	err := baseRequest("GET", fmt.Sprintf("%s/api/v3/calendar?start=%s&end=%s&unmonitored=%v&includeSeries=true", r.InternalAddress, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), unmonitored), nil, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

type getRadarrCalendarEntryResponse struct {
	OriginalTitle   string                         `json:"originalTitle"`
	TitleSlug       string                         `json:"titleSlug"`
	InCinemas       string                         `json:"inCinemas"`
	PhysicalRelease string                         `json:"physicalRelease"`
	DigitalRelease  string                         `json:"digitalRelease"`
	Images          []DefaultReleaseImagesResponse `json:"images"`
	HasFile         bool                           `json:"hasFile"`
}

type DefaultReleaseImagesResponse struct {
	CoverType string `json:"coverType"`
	RemoteURL string `json:"remoteUrl"`
}

func (r *Radarr) GetHealth() ([]*HealthEntry, error) {
	var entries []*HealthEntry
	err := baseRequest("GET", fmt.Sprintf("%s/api/v3/health", r.InternalAddress), nil, &entries)
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
