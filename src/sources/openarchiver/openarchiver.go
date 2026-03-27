package openarchiver

import (
	"fmt"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var (
	BackgroundImgURL = "https://openarchiver.com/logo/logo-sq.svg"
	o                *OpenArchiver
)

type OpenArchiver struct {
	Address          string
	InternalAddress  string
	SuperAPIKey      string
	BackgroundImgURL string
}

func New() (*OpenArchiver, error) {
	if o != nil {
		return o, nil
	}

	address := config.GlobalConfigs.OpenArchiver.Address
	internalAddress := config.GlobalConfigs.OpenArchiver.InternalAddress
	superAPIKey := config.GlobalConfigs.OpenArchiver.SuperAPIKey

	newO := &OpenArchiver{}
	err := newO.Init(address, internalAddress, superAPIKey)
	if err != nil {
		return nil, err
	}

	o = newO

	return o, nil
}

func (l *OpenArchiver) Init(address, internalAddress, superAPIKey string) error {
	if address == "" || superAPIKey == "" {
		return fmt.Errorf("OPENARCHIVER_ADDRESS and OPENARCHIVER_TOKEN variables should be set")
	}

	l.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		l.InternalAddress = l.Address
	} else {
		l.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	l.SuperAPIKey = superAPIKey

	return nil
}
