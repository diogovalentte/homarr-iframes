package pihole

import (
	"fmt"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var (
	p                *Pihole
	BackgroundImgURL = "https://miro.medium.com/v2/resize:fit:657/0*7RBpclLFdUJdwNAK.png"
)

type Pihole struct {
	Address         string
	InternalAddress string
	Token           string
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
	address, internalAddress, APIToken := config.GlobalConfigs.Pihole.Address, config.GlobalConfigs.Pihole.InternalAddress, config.GlobalConfigs.Pihole.Token
	if address == "" || APIToken == "" {
		return fmt.Errorf("PIHOLE_ADDRESS and PIHOLE_TOKEN variables should be set")
	}

	p.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		p.InternalAddress = p.Address
	} else {
		p.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	p.Token = APIToken

	return nil
}

// GetMessages gets the messages that appear in the "Pi-hole diagnostic" page
func (p *Pihole) GetMessages() (*Messages, error) {
	var messages Messages
	err := baseRequest("GET", fmt.Sprintf("%s/admin/api_db.php?messages&auth=%s", p.InternalAddress, p.Token), nil, &messages)
	if err != nil {
		return nil, err
	}

	return &messages, nil
}

type Messages struct {
	Messages []*Message `json:"messages"`
}

type Message struct {
	Blob1     interface{} `json:"blob1"`
	Blob2     interface{} `json:"blob2"`
	Type      string      `json:"type"`
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
}
