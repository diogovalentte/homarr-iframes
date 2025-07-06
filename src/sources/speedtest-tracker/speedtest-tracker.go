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
	token           string
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
	address, internalAddress, token := config.GlobalConfigs.SpeedTestTrackerConfigs.Address, config.GlobalConfigs.SpeedTestTrackerConfigs.InternalAddress, config.GlobalConfigs.SpeedTestTrackerConfigs.Token
	if address == "" || token == "" {
		return fmt.Errorf("SPEEDTEST_TRACKER_ADDRESS and SPEEDTEST_TRACKER_TOKEN variables should be set")
	}

	r.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		r.InternalAddress = r.Address
	} else {
		r.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	r.token = token

	return nil
}

func (s *SpeedTestTracker) GetLatestTest() (*Data, error) {
	var test *TestEntry
	err := s.baseRequest("GET", fmt.Sprintf("%s/api/v1/results/latest", s.InternalAddress), nil, &test)
	if err != nil {
		return nil, err
	}

	return test.Data, nil
}

type TestEntry struct {
	Data    *Data  `json:"data"`
	Message string `json:"message"`
}

type Data struct {
	ID                int     `json:"id"`
	Service           string  `json:"service"`
	Ping              float32 `json:"ping"`
	DownloadBitsHuman string  `json:"download_bits_human"`
	UploadBitsHuman   string  `json:"upload_bits_human"`
	Benchmarks        *struct {
		Download *Benchmarks `json:"download"`
		Upload   *Benchmarks `json:"upload"`
		Ping     *Benchmarks `json:"ping"`
	} `json:"benchmarks"`
	Healthy   bool      `json:"healthy"`
	Status    string    `json:"status"`
	Data      *TestData `json:"data"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type Benchmarks struct {
	Bar    string `json:"bar"`
	Passed bool   `json:"passed"`
	Type   string `json:"type"`
	Value  int    `json:"value"`
	Unit   string `json:"unit"`
}

type TestData struct {
	ISP     string `json:"isp"`
	Message string `json:"message"`
	Level   string `json:"level"`
}
