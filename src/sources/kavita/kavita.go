package kavita

import (
	"fmt"
	"strings"
	"sync"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var (
	k                *Kavita
	BackgroundImgURL = "https://avatars.githubusercontent.com/u/75760308"
)

type Kavita struct {
	Address         string
	InternalAddress string
	Username        string
	Password        string
	Token           string
	RefreshToken    string

	// tokensMutex guards Token and RefreshToken, they're read and written from
	// the goroutines of concurrent requests
	tokensMutex sync.RWMutex
	// refreshMutex serializes the whole refresh/login operation, so concurrent
	// requests getting a 401 at the same time don't refresh on top of each other
	refreshMutex sync.Mutex
}

func (k *Kavita) getTokens() (string, string) {
	k.tokensMutex.RLock()
	defer k.tokensMutex.RUnlock()

	return k.Token, k.RefreshToken
}

func (k *Kavita) setTokens(token, refreshToken string) {
	k.tokensMutex.Lock()
	defer k.tokensMutex.Unlock()

	k.Token, k.RefreshToken = token, refreshToken
}

// newMutex serializes New() calls, otherwise concurrent requests can each
// build their own instance and race on the package singleton
var newMutex sync.Mutex

func New() (*Kavita, error) {
	newMutex.Lock()
	defer newMutex.Unlock()

	if k != nil {
		return k, nil
	}

	newR := &Kavita{}
	err := newR.Init()
	if err != nil {
		return nil, err
	}

	k = newR

	return k, nil
}

func (k *Kavita) Init() error {
	address, internalAddress, username, password := config.GlobalConfigs.Kavita.Address, config.GlobalConfigs.Kavita.InternalAddress, config.GlobalConfigs.Kavita.Username, config.GlobalConfigs.Kavita.Password
	if address == "" || username == "" || password == "" {
		return fmt.Errorf("KAVITA_ADDRESS, KAVITA_USERNAME and KAVITA_PASSWORD variables should be set")
	}

	k.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		k.InternalAddress = k.Address
	} else {
		k.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	k.Username = username
	k.Password = password

	err := k.Login()
	if err != nil {
		return err
	}

	return nil
}

func (k *Kavita) GetMediaErrors() ([]*MediaError, error) {
	var errors MediaErrorResults
	err := k.baseRequest("GET", fmt.Sprintf("%s/api/Server/media-errors", k.InternalAddress), nil, &errors)
	if err != nil {
		return nil, err
	}

	return errors.Results, nil
}

type MediaErrorResults struct {
	Results []*MediaError `json:"results"`
}

type MediaError struct {
	Comment    string `json:"comment"`
	CreatedUTC string `json:"createdUtc"`
}
