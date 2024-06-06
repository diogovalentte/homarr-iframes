package media

import (
	"testing"
	"time"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

func TestGetRadarrCalendar(t *testing.T) {
	s, err := NewRadarr(config.GlobalConfigs.Radarr.Address, config.GlobalConfigs.Radarr.APIKey)
	if err != nil {
		t.Fatalf("error creating Radarr instance: %v", err)
	}
	_, err = s.GetCalendar(false, time.Now(), time.Now().AddDate(0, 0, 1), "inCinemas")
	if err != nil {
		t.Fatalf("error getting calendar: %v", err)
	}
}
