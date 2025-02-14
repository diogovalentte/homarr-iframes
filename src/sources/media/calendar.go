package media

import (
	"fmt"
	"time"

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources/radarr"
	"github.com/diogovalentte/homarr-iframes/src/sources/sonarr"
)

func getCalendar(unmonitored, inCinemas, physical, digital bool) (*Calendar, error) {
	var isAnySourceValid bool
	calendar := &Calendar{}
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 1)

	if config.GlobalConfigs.Radarr.Address != "" && config.GlobalConfigs.Radarr.APIKey != "" {
		isAnySourceValid = true
		radarrCalendar, err := getRadarrCalendar(unmonitored, startDate, endDate, inCinemas, physical, digital)
		if err != nil {
			return nil, fmt.Errorf("couldn't create Radarr calendar: %s", err.Error())
		}
		calendar.Releases = append(calendar.Releases, radarrCalendar.Releases...)
	}

	if config.GlobalConfigs.Sonarr.Address != "" && config.GlobalConfigs.Sonarr.APIKey != "" {
		isAnySourceValid = true
		sonarrCalendar, err := getSonarrCalendar(unmonitored, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("couldn't create Sonarr calendar: %s", err.Error())
		}
		calendar.Releases = append(calendar.Releases, sonarrCalendar.Releases...)
	}

	if !isAnySourceValid {
		return nil, fmt.Errorf("no valid source found. Please check the docs for what environment variables should be set")
	}

	return calendar, nil
}

func getRadarrCalendar(unmonitored bool, startDate, endDate time.Time, inCinemas, physical, digital bool) (*Calendar, error) {
	radarr, err := radarr.New()
	if err != nil {
		return nil, fmt.Errorf("couldn't create Radarr client: %s", err.Error())
	}
	entries, err := radarr.GetCalendar(unmonitored, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("couldn't get Radarr calendar: %s", err.Error())
	}

	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), endDate.Location())

	calendar := &Calendar{}

	for _, entry := range entries {
		var shouldBeDownloaded, found bool
		var releaseDate time.Time

		if inCinemas && entry.InCinemas != "" {
			releaseDate, err = time.Parse(time.RFC3339, entry.InCinemas)
			if err != nil {
				return nil, fmt.Errorf("error parsing movie '%#v' in cinemas date: %w", entry, err)
			}
			releaseDate = releaseDate.In(time.Local)
			if IsReleaseDateWithinDateRange(releaseDate, startDate, endDate) {
				found = true
			}
		}
		if !found && digital && entry.DigitalRelease != "" {
			releaseDate, err := time.Parse(time.RFC3339, entry.DigitalRelease)
			if err != nil {
				return nil, fmt.Errorf("error parsing movie '%#v' digital release date: %w", entry, err)
			}
			releaseDate = releaseDate.In(time.Local)
			if IsReleaseDateWithinDateRange(releaseDate, startDate, endDate) {
				found = true
			}
		}
		if !found && physical && entry.PhysicalRelease != "" {
			releaseDate, err = time.Parse(time.RFC3339, entry.PhysicalRelease)
			if err != nil {
				return nil, fmt.Errorf("error parsing movie '%#v' physical release date: %w", entry, err)
			}
			releaseDate = releaseDate.In(time.Local)
			if IsReleaseDateWithinDateRange(releaseDate, startDate, endDate) {
				found = true
			}
		}
		if !found {
			continue
		}
		if entry.DigitalRelease != "" {
			digitalReleaseDate, err := time.Parse(time.RFC3339, entry.DigitalRelease)
			if err != nil {
				return nil, fmt.Errorf("error parsing movie '%#v' digital release date: %w", entry, err)
			}
			digitalReleaseDate = digitalReleaseDate.In(time.Local)
			shouldBeDownloaded = digitalReleaseDate.Before(time.Now()) && !digitalReleaseDate.IsZero()
		}

		coverImageURL := GetReleaseCoverImageURL(entry.Images)

		calendar.Releases = append(calendar.Releases, MediaRelease{
			Title:              entry.OriginalTitle,
			Source:             "Radarr",
			ReleaseDate:        releaseDate,
			Slug:               entry.TitleSlug,
			CoverImageURL:      coverImageURL,
			IsDownloaded:       entry.HasFile,
			ShouldBeDownloaded: shouldBeDownloaded,
		})
	}

	return calendar, nil
}

func getSonarrCalendar(unmonitored bool, startDate, endDate time.Time) (*Calendar, error) {
	sonarr, err := sonarr.New()
	if err != nil {
		return nil, fmt.Errorf("couldn't create Sonarr client: %s", err.Error())
	}
	entries, err := sonarr.GetCalendar(unmonitored, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("couldn't get Sonarr calendar: %s", err.Error())
	}

	calendar := &Calendar{}

	for _, entry := range entries {
		coverImageURL := GetReleaseCoverImageURL(entry.Series.Images)
		airDate, err := time.Parse(time.RFC3339, entry.AirDateUTC)
		if err != nil {
			return nil, fmt.Errorf("error parsing episode '%#v' air date: %w", entry, err)
		}
		airDate = airDate.In(time.Local)
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
				EpisodeName   string
				SeasonNumber  int
				EpisodeNumber int
			}{
				SeasonNumber:  entry.SeasonNumber,
				EpisodeNumber: entry.EpisodeNumber,
				EpisodeName:   entry.EpisodeTitle,
			},
		})
	}

	return calendar, nil
}
