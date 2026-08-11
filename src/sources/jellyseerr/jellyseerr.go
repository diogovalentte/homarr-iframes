package jellyseerr

import (
	"fmt"
	"strings"
	"sync"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var j *Jellyseerr

type Jellyseerr struct {
	Address         string
	InternalAddress string
	APIKey          string
}

// newMutex serializes New() calls, otherwise concurrent requests can each
// build their own instance and race on the package singleton
var newMutex sync.Mutex

func New() (*Jellyseerr, error) {
	newMutex.Lock()
	defer newMutex.Unlock()

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
