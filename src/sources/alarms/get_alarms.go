package alarms

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/diogovalentte/homarr-iframes/src/config"
	"github.com/diogovalentte/homarr-iframes/src/sources/changedetectionio"
	"github.com/diogovalentte/homarr-iframes/src/sources/kaizoku"
	"github.com/diogovalentte/homarr-iframes/src/sources/kavita"
	"github.com/diogovalentte/homarr-iframes/src/sources/netdata"
	"github.com/diogovalentte/homarr-iframes/src/sources/pihole"
	"github.com/diogovalentte/homarr-iframes/src/sources/prowlarr"
	"github.com/diogovalentte/homarr-iframes/src/sources/radarr"
	"github.com/diogovalentte/homarr-iframes/src/sources/sonarr"
	speedtesttracker "github.com/diogovalentte/homarr-iframes/src/sources/speedtest-tracker"
)

var validAlarmNames = []string{"netdata", "prowlarr", "sonarr", "radarr", "speedtest-tracker", "pihole", "kavita", "kaizoku", "changedetectionio"}

func (a *Alarms) GetAlarms(alarmNames []string, limit int, desc, changedetectionioShowViewed bool) ([]Alarm, error) {
	if limit == 0 {
		return []Alarm{}, nil
	}

	var alarms []Alarm

	for _, alarmName := range alarmNames {
		switch alarmName {
		case "netdata":
			netdataAlarms, err := getNetdataAlarms()
			if err != nil {
				return nil, fmt.Errorf("failed to get Netdata alarms: %w", err)
			}
			alarms = append(alarms, netdataAlarms...)
		case "prowlarr":
			prowlarrAlarms, err := getProwlarrAlarms()
			if err != nil {
				return nil, fmt.Errorf("failed to get Prowlarr alarms: %w", err)
			}
			alarms = append(alarms, prowlarrAlarms...)
		case "radarr":
			radarrAlarms, err := getRadarrAlarms()
			if err != nil {
				return nil, fmt.Errorf("failed to get Radarr alarms: %w", err)
			}
			alarms = append(alarms, radarrAlarms...)
		case "sonarr":
			sonarrAlarms, err := getSonarrAlarms()
			if err != nil {
				return nil, fmt.Errorf("failed to get Sonarr alarms: %w", err)
			}
			alarms = append(alarms, sonarrAlarms...)
		case "speedtest-tracker":
			speedTestTrackerAlarms, err := getSpeedTestTrackerAlarms()
			if err != nil {
				return nil, fmt.Errorf("failed to get SpeedTest Tracker alarms: %w", err)
			}
			alarms = append(alarms, speedTestTrackerAlarms...)
		case "pihole":
			piholeAlarms, err := getPiholeAlarms()
			if err != nil {
				return nil, fmt.Errorf("failed to get Pi-hole alarms: %w", err)
			}
			alarms = append(alarms, piholeAlarms...)
		case "kavita":
			kavitaAlarms, err := getKavitaAlarms()
			if err != nil {
				return nil, fmt.Errorf("failed to get Kavita alarms: %w", err)
			}
			alarms = append(alarms, kavitaAlarms...)
		case "kaizoku":
			kaizokuAlarms, err := getKaizokuAlarms()
			if err != nil {
				return nil, fmt.Errorf("failed to get Kaizoku alarms: %w", err)
			}
			alarms = append(alarms, kaizokuAlarms...)
		case "changedetectionio":
			changedetectionioAlarms, err := getChangedetectionioAlarms(changedetectionioShowViewed)
			if err != nil {
				return nil, fmt.Errorf("failed to get ChangeDetection.io alarms: %w", err)
			}
			alarms = append(alarms, changedetectionioAlarms...)
		default:
			return nil, fmt.Errorf("invalid alarm name: %s", alarmName)
		}
	}

	if len(alarms) > limit && limit > 0 {
		alarms = alarms[:limit]
	}

	sortAlarms(alarms, desc)

	return alarms, nil
}

func getNetdataAlarms() ([]Alarm, error) {
	n, err := netdata.New()
	if err != nil {
		return nil, err
	}

	netdataAlarms, err := n.GetAlarms(-1)
	if err != nil {
		return nil, err
	}

	var alarms []Alarm
	for _, alarm := range netdataAlarms {
		summary := alarm.Summary
		if summary == "" {
			summary = alarm.Name
		}
		if summary == "" {
			summary = "Unknown"
		}
		alarms = append(alarms, Alarm{
			Source:            "Netdata",
			BackgroundImgURL:  netdata.BackgroundImageURL,
			BackgroundImgSize: 80,
			Summary:           summary,
			URL:               n.Address,
			Status:            alarm.Status,
			Value:             alarm.ValueString,
			Property:          alarm.Component + " / " + alarm.Type,
			Time:              alarm.LastStatusChange,
		})
	}

	return alarms, nil
}

func getRadarrAlarms() ([]Alarm, error) {
	r, err := radarr.New()
	if err != nil {
		return nil, err
	}

	radarrAlarms, err := r.GetHealth()
	if err != nil {
		return nil, err
	}

	var alarms []Alarm
	for _, alarm := range radarrAlarms {
		url := fmt.Sprintf("%s/system/status", r.Address)
		if alarm.WikiURL != "" {
			url = alarm.WikiURL
		}
		alarms = append(alarms, Alarm{
			Source:            "Radarr",
			BackgroundImgURL:  radarr.BackgroundImageURL,
			BackgroundImgSize: 120,
			Summary:           alarm.Message,
			URL:               url,
			Status:            strings.ToUpper(alarm.Type),
			Property:          alarm.Source,
		})
	}

	return alarms, nil
}

func getSonarrAlarms() ([]Alarm, error) {
	s, err := sonarr.New()
	if err != nil {
		return nil, err
	}

	sonarrAlarms, err := s.GetHealth()
	if err != nil {
		return nil, err
	}

	var alarms []Alarm
	for _, alarm := range sonarrAlarms {
		url := fmt.Sprintf("%s/system/status", s.Address)
		if alarm.WikiURL != "" {
			url = alarm.WikiURL
		}
		alarms = append(alarms, Alarm{
			Source:            "Sonarr",
			BackgroundImgURL:  sonarr.BackgroundImageURL,
			BackgroundImgSize: 120,
			Summary:           alarm.Message,
			URL:               url,
			Status:            strings.ToUpper(alarm.Type),
			Property:          alarm.Source,
		})
	}

	return alarms, nil
}

func getProwlarrAlarms() ([]Alarm, error) {
	p, err := prowlarr.New()
	if err != nil {
		return nil, err
	}

	prowlarrAlarms, err := p.GetHealth()
	if err != nil {
		return nil, err
	}

	var alarms []Alarm
	for _, alarm := range prowlarrAlarms {
		url := fmt.Sprintf("%s/system/status", p.Address)
		if alarm.WikiURL != "" {
			url = alarm.WikiURL
		}
		alarms = append(alarms, Alarm{
			Source:            "Prowlarr",
			BackgroundImgURL:  prowlarr.BackgroundImageURL,
			BackgroundImgSize: 100,
			Summary:           alarm.Message,
			URL:               url,
			Status:            strings.ToUpper(alarm.Type),
			Property:          alarm.Source,
		})
	}

	return alarms, nil
}

func getSpeedTestTrackerAlarms() ([]Alarm, error) {
	s, err := speedtesttracker.New()
	if err != nil {
		return nil, err
	}

	test, err := s.GetLatestTest()
	if err != nil {
		return nil, err
	}

	// test failed
	if test.Status == "failed" {
		layout := "2006-01-02 15:04:05"
		updatedAt, err := time.Parse(layout, test.UpdatedAt)
		if err != nil {
			return nil, err
		}

		url := s.Address + "/admin/results"

		alarms := []Alarm{{
			Time:            updatedAt,
			Summary:         "Last Speedtest Failed",
			URL:             url,
			Status:          strings.ToUpper(test.Data.Level),
			Value:           test.Service,
			Property:        test.Data.Message,
			Source:          "SpeedTest",
			BackgroundColor: "black",
		}}

		return alarms, nil
	} else if test.Status == "running" {
		return []Alarm{}, nil
	}

	// threshold breached
	if !test.Healthy {
		layout := "2006-01-02 15:04:05"
		updatedAt, err := time.Parse(layout, test.UpdatedAt)
		if err != nil {
			return nil, err
		}

		url := s.Address + "/admin/results"

		var status, value string
		if !test.Benchmarks.Download.Passed {
			status = "DOWNLOAD"
			value = test.DownloadBitsHuman
		} else if !test.Benchmarks.Upload.Passed {
			status = "UPLOAD"
			value = test.UploadBitsHuman
		} else if !test.Benchmarks.Ping.Passed {
			status = "PING"
			value = fmt.Sprintf("%s ms", test.Ping)
		}

		alarms := []Alarm{{
			Time:            updatedAt,
			Summary:         "Speedtest Threshold Breached",
			URL:             url,
			Status:          status,
			Value:           value,
			Property:        test.Data.ISP,
			Source:          "SpeedTest",
			BackgroundColor: "black",
		}}

		return alarms, nil
	}

	return []Alarm{}, nil
}

func getPiholeAlarms() ([]Alarm, error) {
	p, err := pihole.New()
	if err != nil {
		return nil, err
	}

	messages, err := p.GetMessages()
	if err != nil {
		return nil, err
	}

	var alarms []Alarm
	for _, message := range messages.Messages {
		timestamp := time.Unix(message.Timestamp, 0)
		var url string
		if p.Token != "" {
			url = p.Address + "/admin/messages.php"
		} else {
			url = p.Address + "/admin/messages"
		}
		alarms = append(alarms, Alarm{
			Source:            "Pi-hole",
			BackgroundImgURL:  pihole.BackgroundImgURL,
			BackgroundImgSize: 80,
			Summary:           message.Plain,
			Property:          message.Type,
			Time:              timestamp,
			URL:               url,
			Status:            "WARNING",
		})
	}

	return alarms, nil
}

func getKavitaAlarms() ([]Alarm, error) {
	k, err := kavita.New()
	if err != nil {
		return nil, err
	}

	errors, err := k.GetMediaErrors()
	if err != nil {
		return nil, err
	}

	var alarms []Alarm
	for _, error := range errors {
		layout := "2006-01-02T15:04:05.9999999"
		timestamp, err := time.Parse(layout, error.CreatedUTC)
		if err != nil {
			return nil, err
		}
		timestamp = timestamp.Local()
		url := k.Address + "/settings#admin-media-issues"
		alarms = append(alarms, Alarm{
			Source:            "Kavita",
			BackgroundImgURL:  kavita.BackgroundImgURL,
			BackgroundImgSize: 100,
			Summary:           error.Comment,
			Time:              timestamp,
			URL:               url,
			Status:            "ERROR",
		})
	}

	return alarms, nil
}

func getKaizokuAlarms() ([]Alarm, error) {
	k, err := kaizoku.New()
	if err != nil {
		return nil, err
	}

	queues, err := k.GetQueues()
	if err != nil {
		return nil, err
	}

	var alarms []Alarm
	for _, queue := range queues {
		if queue.Counts.Failed < 1 {
			continue
		}

		summary := fmt.Sprintf("Queue %s has failed jobs", queue.Name)
		url := k.Address + "/bull/queues/queue/" + queue.Name + "?status=failed"

		alarms = append(alarms, Alarm{
			Source:          "Kaizoku",
			BackgroundColor: "black",
			Summary:         summary,
			URL:             url,
			Status:          "FAILED",
			Value:           fmt.Sprintf("%d jobs", queue.Counts.Failed),
			Property:        queue.Name,
		})
	}

	return alarms, nil
}

func getChangedetectionioAlarms(showViewed bool) ([]Alarm, error) {
	c, err := changedetectionio.New()
	if err != nil {
		return nil, err
	}

	watches, err := c.GetWatches()
	if err != nil {
		return nil, err
	}

	var alarms []Alarm
	for ID, watch := range watches {
		var hasError bool
		if value, ok := watch.LastError.(bool); !ok {
			errStr := watch.LastError.(string)
			if errStr != "" {
				hasError = true
			}
		} else if value {
			hasError = true
		}

		if hasError {
			errStr := watch.LastError.(string)
			lastChecked := time.Unix(int64(watch.LastChecked), 0)
			summary := watch.Title
			if watch.Title == "" {
				summary = watch.URL
			}
			alarms = append(alarms, Alarm{
				Source:            "ChangeDetection.io",
				BackgroundImgURL:  changedetectionio.BackgroundImgURL,
				BackgroundImgSize: 100,
				Summary:           summary,
				URL:               c.Address,
				Status:            "ERROR",
				Property:          errStr,
				Time:              lastChecked,
			})

			continue
		}

		if watch.Viewed && !showViewed {
			continue
		}

		minChanged := time.Now().Add(-time.Hour * time.Duration(config.GlobalConfigs.ChangeDetectionIO.ChangedLastHours))
		lastChanged := time.Unix(int64(watch.LastChanged), 0)
		if lastChanged.After(minChanged) {
			summary := watch.Title
			if summary == "" {
				summary = watch.URL
			}
			viewed := "Viewed"
			if !watch.Viewed {
				viewed = "Not Viewed"
			}

			alarms = append(alarms, Alarm{
				Source:            "ChangeDetection.io",
				BackgroundImgURL:  changedetectionio.BackgroundImgURL,
				BackgroundImgSize: 100,
				Summary:           summary,
				URL:               c.Address + "/diff/" + ID,
				Status:            "CHANGED",
				Value:             viewed,
				Time:              lastChanged,
			})
		}
	}

	return alarms, nil
}

func sortAlarms(alarms []Alarm, desc bool) {
	if desc {
		sort.Slice(alarms, func(i, j int) bool {
			if alarms[i].Time.IsZero() && alarms[j].Time.IsZero() {
				return false
			}
			if alarms[i].Time.IsZero() {
				return false
			}
			if alarms[j].Time.IsZero() {
				return true
			}
			return alarms[i].Time.After(alarms[j].Time)
		})
		return
	}
	sort.Slice(alarms, func(i, j int) bool {
		if alarms[i].Time.IsZero() && alarms[j].Time.IsZero() {
			return false
		}
		if alarms[i].Time.IsZero() {
			return false
		}
		if alarms[j].Time.IsZero() {
			return true
		}
		return alarms[i].Time.Before(alarms[j].Time)
	})
}
