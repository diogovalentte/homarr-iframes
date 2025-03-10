package pihole

import (
	"fmt"
	"strings"
	"time"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var (
	p                *Pihole
	BackgroundImgURL = "https://miro.medium.com/v2/resize:fit:657/0*7RBpclLFdUJdwNAK.png"
)

type Pihole struct {
	Address         string
	InternalAddress string
	Token           string // <v6.0
	SID             string
	Password        string
	ValidityTime    time.Time
}

func New() (*Pihole, error) {
	if p != nil {
		return p, nil
	}

	newR := &Pihole{}
	err := newR.Init()
	if err != nil {
		return nil, err
	}

	p = newR

	return p, nil
}

func (p *Pihole) Init() error {
	address, internalAddress, APIToken, APIPassword := config.GlobalConfigs.Pihole.Address, config.GlobalConfigs.Pihole.InternalAddress, config.GlobalConfigs.Pihole.Token, config.GlobalConfigs.Pihole.Password
	if address == "" || (APIToken == "" && APIPassword == "") {
		return fmt.Errorf("PIHOLE_ADDRESS and PIHOLE_TOKEN or PIHOLE_PASSWORD variables should be set")
	}

	p.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		p.InternalAddress = p.Address
	} else {
		p.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}

	if APIToken != "" {
		p.Token = APIToken
	} else {
		p.Password = APIPassword
		err := p.Login()
		if err != nil {
			return err
		}
	}

	return nil
}

// GetMessages gets the messages that appear in the "Pi-hole diagnostic" page
func (p *Pihole) GetMessages() (*Messages, error) {
	var messages Messages
	var url string
	if p.Token != "" {
		url = fmt.Sprintf("%s/admin/api.php?messages?auth=%s", p.InternalAddress, p.Token)
	} else {
		url = fmt.Sprintf("%s/api/info/messages", p.InternalAddress)
	}
	err := p.baseRequest("GET", url, nil, &messages, 1)
	if err != nil {
		return nil, err
	}

	return &messages, nil
}

type Messages struct {
	Messages []*Message `json:"messages"`
}

type Message struct {
	Type      string `json:"type"`
	Plain     string `json:"plain"`
	HTML      string `json:"html"`
	Timestamp int64  `json:"timestamp"`
}
