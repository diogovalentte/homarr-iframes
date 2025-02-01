package changedetectionio

import (
	"fmt"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var (
	c                *ChangeDetectionIO
	BackgroundImgURL = "https://i.imgur.com/16Q6GPD.png"
)

type ChangeDetectionIO struct {
	Address         string
	InternalAddress string
	APIKey          string
}

func New() (*ChangeDetectionIO, error) {
	if c != nil {
		return c, nil
	}

	newR := &ChangeDetectionIO{}
	err := newR.Init()
	if err != nil {
		return nil, err
	}

	c = newR

	return c, nil
}

func (c *ChangeDetectionIO) Init() error {
	address, internalAddress, APIKey := config.GlobalConfigs.ChangeDetectionIO.Address, config.GlobalConfigs.ChangeDetectionIO.InternalAddress, config.GlobalConfigs.ChangeDetectionIO.APIKey
	if address == "" || APIKey == "" {
		return fmt.Errorf("CHANGEDETECTIONIO_ADDRESS and CHANGEDETECTIONIO_API_KEY variables should be set")
	}

	c.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		c.InternalAddress = c.Address
	} else {
		c.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	c.APIKey = APIKey

	return nil
}

func (c *ChangeDetectionIO) GetWatches() (map[string]*Watch, error) {
	var watches map[string]*Watch
	err := c.baseRequest("GET", fmt.Sprintf("%s/api/v1/watch", c.Address), nil, &watches)
	if err != nil {
		return nil, err
	}

	return watches, nil
}

type Watch struct {
	Title       string      `json:"title"`
	URL         string      `json:"url"`
	LastChanged int         `json:"last_changed"`
	LastChecked int         `json:"last_checked"`
	Viewed      bool        `json:"viewed"`
	LastError   interface{} `json:"last_error"`
}
