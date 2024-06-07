package media

import (
	"fmt"
	"strings"
	"time"
)

var r *Radarr

type Radarr struct {
	Address string
	APIKey  string
}

func NewRadarr(address, APIKey string) (*Radarr, error) {
	if r != nil {
		return r, nil
	}

	newR := &Radarr{}
	err := newR.Init(address, APIKey)
	if err != nil {
		return nil, err
	}

	r = newR

	return r, nil
}

func (r *Radarr) Init(address, token string) error {
	if address == "" || token == "" {
		return fmt.Errorf("RADARR_ADDRESS and RADARR_API_KEY variables should be set")
	}

	if strings.HasSuffix(address, "/") {
		address = address[:len(address)-1]
	}

	r.Address = address
	r.APIKey = token

	return nil
}

// GetCalendar returns the calendar of releases where the release date is between startDate and endDate.
// To get the calendar of a specific day, set the startDate to the specific day and endDate to one day after.
// To get the calendar of a specific week, set the startDate to the first day of the week and endDate to one day after the last day of the week.
// It considers only the date, not the time.
func (r *Radarr) GetCalendar(unmonitored bool, startDate, endDate time.Time, releaseType string) (*Calendar, error) {
	calendar := &Calendar{}
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	var entries []getRadarrCalendarEntryResponse
	err := baseRequest("GET", fmt.Sprintf("%s/api/v3/calendar?apiKey=%s&start=%s&end=%s&unmonitored=%v&includeSeries=true", r.Address, r.APIKey, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), unmonitored), nil, &entries)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		var releaseDate time.Time
		switch releaseType {
		case "inCinemas":
			if entry.InCinemas == "" {
				continue
			}
			releaseDate, err = time.Parse(time.RFC3339, entry.InCinemas)
			if err != nil {
				return nil, fmt.Errorf("error parsing movie '%#v' in cinemas date: %w", entry, err)
			}
			releaseDate = releaseDate.In(time.Local)
			if !isReleaseDateWithinDateRange(releaseDate, startDate, endDate) {
				continue
			}
		case "physical":
			if entry.PhysicalRelease == "" {
				continue
			}
			releaseDate, err = time.Parse(time.RFC3339, entry.InCinemas)
			if err != nil {
				return nil, fmt.Errorf("error parsing movie '%#v' physical release date: %w", entry, err)
			}
			releaseDate = releaseDate.In(time.Local)
			if !isReleaseDateWithinDateRange(releaseDate, startDate, endDate) {
				continue
			}
		case "digital":
			if entry.DigitalRelease == "" {
				continue
			}
			releaseDate, err = time.Parse(time.RFC3339, entry.InCinemas)
			if err != nil {
				return nil, fmt.Errorf("error parsing movie '%#v' digital release date: %w", entry, err)
			}
			releaseDate = releaseDate.In(time.Local)
			if !isReleaseDateWithinDateRange(releaseDate, startDate, endDate) {
				continue
			}
		default:
			return nil, fmt.Errorf("invalid release type: %s", releaseType)
		}

		coverImageURL := getReleaseCoverImageURL(entry.Images)

		calendar.Releases = append(calendar.Releases, MediaRelease{
			Title:         entry.OriginalTitle,
			Source:        "Radarr",
			ReleaseDate:   releaseDate,
			Slug:          entry.TitleSlug,
			CoverImageURL: coverImageURL,
			IsDownloaded:  entry.HasFile,
		})
	}

	return calendar, nil
}

type getRadarrCalendarEntryResponse struct {
	OriginalTitle   string                         `json:"originalTitle"`
	HasFile         bool                           `json:"hasFile"`
	TitleSlug       string                         `json:"titleSlug"`
	InCinemas       string                         `json:"inCinemas"`
	PhysicalRelease string                         `json:"physicalRelease"`
	DigitalRelease  string                         `json:"digitalRelease"`
	Images          []defaultReleaseImagesResponse `json:"images"`
}
