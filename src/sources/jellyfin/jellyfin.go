package jellyfin

import (
	"fmt"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

type Jellyfin struct {
	Address         string
	InternalAddress string
	APIKey          string
	userId          string
}

var j *Jellyfin

func New() (*Jellyfin, error) {
	if j != nil {
		return j, nil
	}

	address := config.GlobalConfigs.Jellyfin.Address
	internalAddress := config.GlobalConfigs.Jellyfin.InternalAddress
	APIKey := config.GlobalConfigs.Jellyfin.APIKey
	UserId := config.GlobalConfigs.Jellyfin.UserId

	newj := &Jellyfin{}
	err := newj.Init(address, internalAddress, APIKey, UserId)
	if err != nil {
		return nil, err
	}

	j = newj

	return j, nil
}

func (j *Jellyfin) Init(address, internalAddress, APIKey, UserId string) error {
	if address == "" || APIKey == "" {
		return fmt.Errorf("JELLYFIN_ADDRESS and JELLYFIN_API_KEY variables should be set")
	}

	j.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		j.InternalAddress = j.Address
	} else {
		j.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	j.APIKey = APIKey
	j.userId = UserId

	return nil
}
