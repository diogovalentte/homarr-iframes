package netdata

import (
	"fmt"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var (
	n                  *Netdata
	BackgroundImageURL = "https://avatars.githubusercontent.com/u/43390781"
)

type Netdata struct {
	Address         string
	InternalAddress string
	Token           string
}

func New() (*Netdata, error) {
	if n != nil {
		return n, nil
	}

	newN := &Netdata{}
	err := newN.Init()
	if err != nil {
		return nil, err
	}

	n = newN

	return n, nil
}

// Init sets the Netdata properties from the configs
func (n *Netdata) Init() error {
	address, internalAddress, token := config.GlobalConfigs.NetdataConfigs.Address, config.GlobalConfigs.NetdataConfigs.InternalAddress, config.GlobalConfigs.NetdataConfigs.Token
	if address == "" || token == "" {
		return fmt.Errorf("NETDATA_ADDRESS and NETDATA_TOKEN variables should be set")
	}
	n.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		n.InternalAddress = n.Address
	} else {
		n.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	n.Token = token

	return nil
}
