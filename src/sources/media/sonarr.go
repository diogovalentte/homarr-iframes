package media

import (
	"fmt"
	"strings"
	"time"
)

var s *Sonarr

type Sonarr struct {
	Address         string
	InternalAddress string
	APIKey          string
}

func NewSonarr(address, internalAddress, APIKey string) (*Sonarr, error) {
	if s != nil {
		return s, nil
	}

	newS := &Sonarr{}
	err := newS.Init(address, internalAddress, APIKey)
	if err != nil {
		return nil, err
	}

	s = newS

	return s, nil
}

func (s *Sonarr) Init(address, internalAddress, APIKey string) error {
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

// GetCalendar returns the calendar of releases where the air date is between startDate and endDate.
// To get the calendar of a specific day, set the startDate to the specific day and endDate to one day after.
// To get the calendar of a specific week, set the startDate to the first day of the week and endDate to one day after the last day of the week.
// It considers only the date, not the time.
func (s *Sonarr) GetCalendar(unmonitored bool, startDate, endDate time.Time) (*Calendar, error) {
	calendar := &Calendar{}
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	var entries []getSonarrCalendarEntryResponse
	err := baseRequest("GET", fmt.Sprintf("%s/api/v3/calendar?apiKey=%s&start=%s&end=%s&unmonitored=%v&includeSeries=true", s.InternalAddress, s.APIKey, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), unmonitored), nil, &entries)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		coverImageURL := getReleaseCoverImageURL(entry.Series.Images)
		airDate, err := time.Parse(time.RFC3339, entry.AirDateUTC)
		if err != nil {
			return nil, fmt.Errorf("error parsing episode '%#v' air date: %w", entry, err)
		}
		airDate = airDate.In(time.Local)
		if !isReleaseDateWithinDateRange(airDate, startDate, endDate) {
			continue
		}
		now := time.Now()
		shouldBeDownloaded := airDate.Before(now)

		calendar.Releases = append(calendar.Releases, MediaRelease{
			Title:              entry.Series.Title,
			Source:             "Sonarr",
			ReleaseDate:        airDate,
			Slug:               entry.Series.TitleSlug,
			CoverImageURL:      coverImageURL,
			IsDownloaded:       entry.HasFile,
			ShouldBeDownloaded: shouldBeDownloaded,
			EpisodeDetails: struct {
				SeasonNumber  int
				EpisodeNumber int
				EpisodeName   string
			}{
				SeasonNumber:  entry.SeasonNumber,
				EpisodeNumber: entry.EpisodeNumber,
				EpisodeName:   entry.EpisodeTitle,
			},
		})
	}

	return calendar, nil
}

type getSonarrCalendarEntryResponse struct {
	SeasonNumber  int    `json:"seasonNumber"`
	EpisodeNumber int    `json:"episodeNumber"`
	EpisodeTitle  string `json:"title"`
	HasFile       bool   `json:"hasFile"`
	AirDateUTC    string `json:"airDateUtc"`
	Series        struct {
		Title     string                         `json:"title"`
		TitleSlug string                         `json:"titleSlug"`
		Images    []defaultReleaseImagesResponse `json:"images"`
	} `json:"series"`
}
