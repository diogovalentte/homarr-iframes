package jellyseerr

import (
	"fmt"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var (
	j                  *Jellyseerr
	tmdbImageBasePath  = "https://image.tmdb.org/t/p/w600_and_h900_bestv2/"
	backgroundImageURL = "https://github.com/Fallenbagel/jellyseerr/blob/develop/public/logo_full.png"
)

type Jellyseerr struct {
	Address         string
	InternalAddress string
	APIKey          string
}

func New() (*Jellyseerr, error) {
	if j != nil {
		return j, nil
	}

	address := config.GlobalConfigs.Jellyseerr.Address
	internalAddress := config.GlobalConfigs.Jellyseerr.InternalAddress
	APIKey := config.GlobalConfigs.Jellyseerr.APIKey

	newJ := &Jellyseerr{}
	err := newJ.Init(address, internalAddress, APIKey)
	if err != nil {
		return nil, err
	}

	j = newJ

	return j, nil
}

// Init sets the jellyseerr properties from the configs
func (j *Jellyseerr) Init(address, internalAddress, APIKey string) error {
	if address == "" || APIKey == "" {
		return fmt.Errorf("JELLYSEERR_ADDRESS and JELLYSEERR_API_KEY variables should be set")
	}

	j.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		j.InternalAddress = j.Address
	} else {
		j.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	j.APIKey = APIKey

	return nil
}
