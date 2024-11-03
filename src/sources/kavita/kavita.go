package kavita

import (
	"fmt"
	"strings"

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
}

func New() (*Kavita, error) {
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
	var errors []*MediaError
	err := k.baseRequest("GET", fmt.Sprintf("%s/api/Server/media-errors", k.InternalAddress), nil, &errors)
	if err != nil {
		return nil, err
	}

	return errors, nil
}

type MediaError struct {
	Comment    string `json:"comment"`
	CreatedUTC string `json:"createdUtc"`
}
