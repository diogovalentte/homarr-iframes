package sonarr

import (
	"fmt"
	"strings"
	"time"

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources/radarr"
)

var (
	s                  *Sonarr
	BackgroundImageURL = "https://avatars.githubusercontent.com/u/1082903"
)

type Sonarr struct {
	Address         string
	InternalAddress string
	APIKey          string
}

func New() (*Sonarr, error) {
	if s != nil {
		return s, nil
	}

	newS := &Sonarr{}
	err := newS.Init()
	if err != nil {
		return nil, err
	}

	s = newS

	return s, nil
}

func (s *Sonarr) Init() error {
	address, internalAddress, APIKey := config.GlobalConfigs.Sonarr.Address, config.GlobalConfigs.Sonarr.InternalAddress, config.GlobalConfigs.Sonarr.APIKey
	if address == "" || APIKey == "" {
		return fmt.Errorf("SONARR_ADDRESS and SONARR_API_KEY variables should be set")
	}

	s.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		s.InternalAddress = s.Address
	} else {
		s.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	s.APIKey = APIKey

	return nil
}

// GetCalendar returns the calendar of releases where the release date is between startDate and endDate.
// To get the calendar of a specific day, set the startDate and endDate to the same day.
// To get the calendar of a specific week, set the startDate to the first day of the week and endDate to the last day of the week.
// It considers only the date, not the time, so it'll get all releases that are released on that day.
func (s *Sonarr) GetCalendar(unmonitored bool, startDate, endDate time.Time) ([]*getSonarrCalendarEntryResponse, error) {
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), endDate.Location())

	var entries []*getSonarrCalendarEntryResponse
	err := baseRequest("GET", fmt.Sprintf("%s/api/v3/calendar?start=%s&end=%s&unmonitored=%v&includeSeries=true", s.InternalAddress, startDate.Format("2006-01-02T15:04:05.000Z"), endDate.Format("2006-01-02T15:04:05.000Z"), unmonitored), nil, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

type getSonarrCalendarEntryResponse struct {
	EpisodeTitle string `json:"title"`
	AirDateUTC   string `json:"airDateUtc"`
	Series       struct {
		Title     string                                `json:"title"`
		TitleSlug string                                `json:"titleSlug"`
		Images    []radarr.DefaultReleaseImagesResponse `json:"images"`
	} `json:"series"`
	HasFile       bool `json:"hasFile"`
	SeasonNumber  int  `json:"seasonNumber"`
	EpisodeNumber int  `json:"episodeNumber"`
}

func (s *Sonarr) GetHealth() ([]*HealthEntry, error) {
	var entries []*HealthEntry
	err := baseRequest("GET", fmt.Sprintf("%s/api/v3/health", s.InternalAddress), nil, &entries)
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
