package overseerr

import (
	"fmt"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var (
	o                         *Overseerr
	TMDBPosterImageBasePath   = "https://image.tmdb.org/t/p/w600_and_h900_bestv2/"
	TMDBBackdropImageBasePath = "https://image.tmdb.org/t/p/original/"
	DefaultBackgroundImageURL = "https://i.imgur.com/jMy7evE.jpeg"
)

type Overseerr struct {
	Address         string
	InternalAddress string
	APIKey          string
}

func New() (*Overseerr, error) {
	if o != nil {
		return o, nil
	}

	address := config.GlobalConfigs.Overseerr.Address
	internalAddress := config.GlobalConfigs.Overseerr.InternalAddress
	APIKey := config.GlobalConfigs.Overseerr.APIKey

	newO := &Overseerr{}
	err := newO.Init(address, internalAddress, APIKey)
	if err != nil {
		return nil, err
	}

	o = newO

	return o, nil
}

// Init sets the Overseerr properties from the configs
func (o *Overseerr) Init(address, internalAddress, APIKey string) error {
	if address == "" || APIKey == "" {
		return fmt.Errorf("OVERSEERR_ADDRESS and OVERSEERR_API_KEY variables should be set")
	}

	o.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		o.InternalAddress = o.Address
	} else {
		o.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	o.APIKey = APIKey

	return nil
}
