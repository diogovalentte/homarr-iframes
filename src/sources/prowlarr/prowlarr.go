package prowlarr

import (
	"fmt"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var (
	p                  *Prowlarr
	BackgroundImageURL = "https://avatars.githubusercontent.com/u/73049443"
)

type Prowlarr struct {
	Address         string
	InternalAddress string
	APIKey          string
}

func New() (*Prowlarr, error) {
	if p != nil {
		return p, nil
	}

	newR := &Prowlarr{}
	err := newR.Init()
	if err != nil {
		return nil, err
	}

	p = newR

	return p, nil
}

func (p *Prowlarr) Init() error {
	address, internalAddress, APIKey := config.GlobalConfigs.Prowlarr.Address, config.GlobalConfigs.Prowlarr.InternalAddress, config.GlobalConfigs.Prowlarr.APIKey
	if address == "" || APIKey == "" {
		return fmt.Errorf("PROWLARR_ADDRESS and PROWLARR_API_KEY variables should be set")
	}

	p.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		p.InternalAddress = p.Address
	} else {
		p.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	p.APIKey = APIKey

	return nil
}

func (p *Prowlarr) GetHealth() ([]*HealthEntry, error) {
	var entries []*HealthEntry
	err := baseRequest("GET", fmt.Sprintf("%s/api/v1/health", p.InternalAddress), nil, &entries)
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
