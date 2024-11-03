package alarms

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/diogovalentte/homarr-iframes/src/sources/netdata"
	"github.com/diogovalentte/homarr-iframes/src/sources/prowlarr"
	"github.com/diogovalentte/homarr-iframes/src/sources/radarr"
	"github.com/diogovalentte/homarr-iframes/src/sources/sonarr"
	speedtesttracker "github.com/diogovalentte/homarr-iframes/src/sources/speedtest-tracker"
)

var validAlarmNames = []string{"netdata", "prowlarr", "sonarr", "radarr", "speedtest-tracker"}

func (a *Alarms) GetAlarms(alarmNames []string, limit int, desc bool) ([]Alarm, error) {
	if limit == 0 {
		return []Alarm{}, nil
	}

	var alarms []Alarm

	for _, alarmName := range alarmNames {
		switch alarmName {
		case "netdata":
			netdataAlarms, err := getNetdataAlarms()
			if err != nil {
				return nil, err
			}
			alarms = append(alarms, netdataAlarms...)
		case "prowlarr":
			prowlarrAlarms, err := getProwlarrAlarms()
			if err != nil {
				return nil, err
			}
			alarms = append(alarms, prowlarrAlarms...)
		case "radarr":
			radarrAlarms, err := getRadarrAlarms()
			if err != nil {
				return nil, err
			}
			alarms = append(alarms, radarrAlarms...)
		case "sonarr":
			sonarrAlarms, err := getSonarrAlarms()
			if err != nil {
				return nil, err
			}
			alarms = append(alarms, sonarrAlarms...)
		case "speedtest-tracker":
			speedTestTrackerAlarms, err := getSpeedTestTrackerAlarms()
			if err != nil {
				return nil, err
			}
			alarms = append(alarms, speedTestTrackerAlarms...)
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

	var alarms []Alarm

	if !test.Data.Failed {
		return alarms, nil
	}

	url := test.Data.URL
	if url == "" {
		url = s.Address + "/admin/results"
	}
	layout := "2006-01-02T15:04:05.999999999-07:00"
	updatedAt, err := time.Parse(layout, test.Data.UpdatedAt)
	if err != nil {
		return nil, err
	}
	alarms = append(alarms, Alarm{
		Source:          "SpeedTest",
		BackgroundColor: "black",
		Summary:         "Last Speedtest failed",
		URL:             url,
		Status:          "FAILED",
		Property:        test.Data.ServerName,
		Time:            updatedAt,
	})

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
