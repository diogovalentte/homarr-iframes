package speedtesttracker

import (
	"fmt"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var s *SpeedTestTracker

type SpeedTestTracker struct {
	Address         string
	InternalAddress string
}

func New() (*SpeedTestTracker, error) {
	if s != nil {
		return s, nil
	}

	newR := &SpeedTestTracker{}
	err := newR.Init()
	if err != nil {
		return nil, err
	}

	s = newR

	return s, nil
}

func (r *SpeedTestTracker) Init() error {
	address, internalAddress := config.GlobalConfigs.SpeedTestTrackerConfigs.Address, config.GlobalConfigs.SpeedTestTrackerConfigs.InternalAddress
	if address == "" {
		return fmt.Errorf("SPEEDTEST_TRACKER_ADDRESS variable should be set")
	}

	r.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		r.InternalAddress = r.Address
	} else {
		r.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}

	return nil
}

func (s *SpeedTestTracker) GetLatestTest() (*TestEntry, error) {
	var test *TestEntry
	err := baseRequest("GET", fmt.Sprintf("%s/api/speedtest/latest", s.InternalAddress), nil, &test)
	if err != nil {
		return nil, err
	}

	return test, nil
}

type TestEntry struct {
	Data    *TestData `json:"data"`
	Message string    `json:"message"`
}

type TestData struct {
	URL        string `json:"url"`
	UpdatedAt  string `json:"updated_at"`
	ServerName string `json:"server_name"`
	Failed     bool   `json:"failed"`
}
