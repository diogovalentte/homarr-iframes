package media

import (
	"fmt"
	"time"

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources/lidarr"
	"github.com/diogovalentte/homarr-iframes/src/sources/radarr"
	"github.com/diogovalentte/homarr-iframes/src/sources/sonarr"
)

func getCalendar(unmonitored, inCinemas, physical, digital bool) (*Calendar, error) {
	var isAnySourceValid bool
	calendar := &Calendar{}
	startDate := time.Now()
	endDate := startDate

	if config.GlobalConfigs.Radarr.Address != "" && config.GlobalConfigs.Radarr.APIKey != "" {
		isAnySourceValid = true
		radarrCalendar, err := getRadarrCalendar(unmonitored, startDate, endDate, inCinemas, physical, digital)
		if err != nil {
			return nil, fmt.Errorf("couldn't create Radarr calendar: %s", err.Error())
		}
		calendar.Releases = append(calendar.Releases, radarrCalendar.Releases...)
	}

	if config.GlobalConfigs.Lidarr.Address != "" && config.GlobalConfigs.Lidarr.APIKey != "" {
		isAnySourceValid = true
		lidarrCalendar, err := getLidarrCalendar(unmonitored, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("couldn't create Lidarr calendar: %s", err.Error())
		}
		calendar.Releases = append(calendar.Releases, lidarrCalendar.Releases...)
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

		posterImageURL, coverImageURL := GetReleaseImagesURL(entry.Images)

		calendar.Releases = append(calendar.Releases, MediaRelease{
			Title:              entry.OriginalTitle,
			Source:             "Radarr",
			ReleaseDate:        releaseDate,
			Slug:               entry.TitleSlug,
			CoverImageURL:      coverImageURL,
			PosterImageURL:     posterImageURL,
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
		posterImageURL, coverImageURL := GetReleaseImagesURL(entry.Series.Images)
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
			PosterImageURL:     posterImageURL,
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

func getLidarrCalendar(unmonitored bool, startDate, endDate time.Time) (*Calendar, error) {
	lidarr, err := lidarr.New()
	if err != nil {
		return nil, fmt.Errorf("couldn't create Lidarr client: %s", err.Error())
	}
	entries, err := lidarr.GetCalendar(unmonitored, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("couldn't get Lidarr calendar: %s", err.Error())
	}

	calendar := &Calendar{}

	for _, entry := range entries {
		posterImageURL, coverImageURL := GetReleaseImagesURL(entry.Images)
		if posterImageURL == "" {
			posterImageURL, coverImageURL = GetReleaseImagesURL(entry.Artist.Images)
		}
		airDate, err := time.Parse(time.RFC3339, entry.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("error parsing album '%#v' release date: %w", entry, err)
		}
		airDate = airDate.In(time.Local)

		calendar.Releases = append(calendar.Releases, MediaRelease{
			Title:          entry.Title,
			Source:         "Lidarr",
			ReleaseDate:    airDate,
			Slug:           entry.ForeignAlbumID,
			CoverImageURL:  coverImageURL,
			PosterImageURL: posterImageURL,
			IsDownloaded:   false,
			ArtistDetails: struct {
				ArtistName string
				Slug       string
			}{
				ArtistName: entry.Artist.ArtistName,
				Slug:       entry.Artist.ForeignArtistID,
			},
			AlbumType:       entry.AlbumType,
			TrackFileCount:  entry.Statistics.TrackFileCount,
			TotalTrackCount: entry.Statistics.TotalTrackCount,
		})
	}

	return calendar, nil
}

func GetReleaseImagesURL(images []radarr.DefaultReleaseImagesResponse) (string, string) {
	if len(images) == 0 {
		return "", ""
	}

	var posterURL, coverURL string
	for _, image := range images {
		if image.CoverType == "poster" {
			if image.RemoteURL != "" {
				posterURL = image.RemoteURL
			} else {
				posterURL = image.URL
			}
		} else if image.CoverType == "banner" {
			if image.RemoteURL != "" {
				coverURL = image.RemoteURL
			} else {
				coverURL = image.URL
			}
		} else if image.CoverType == "cover" {
			if image.RemoteURL != "" {
				coverURL = image.RemoteURL
			} else {
				coverURL = image.URL
			}
		}
	}

	if posterURL == "" && coverURL == "" {
		if images[0].RemoteURL != "" {
			posterURL, coverURL = images[0].RemoteURL, images[0].RemoteURL
		}
	} else if coverURL == "" {
		coverURL = posterURL
	} else if posterURL == "" {
		posterURL = coverURL
	}

	return posterURL, coverURL
}

// IsReleaseDateWithinDateRange checks if it's within a given date range.
// startDate is inclusive, endDate is exclusive.
func IsReleaseDateWithinDateRange(releaseDate, startDate, endDate time.Time) bool {
	return (releaseDate.After(startDate) || releaseDate.Equal(startDate)) && releaseDate.Before(endDate)
}
