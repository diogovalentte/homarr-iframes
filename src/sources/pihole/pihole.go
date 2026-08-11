package pihole

import (
	"fmt"
	"strings"
	"sync"
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

	// sessionMutex guards SID and ValidityTime, they're read and written from
	// the goroutines of concurrent requests
	sessionMutex sync.RWMutex
	// renewMutex serializes the logout/login pair, otherwise a request can log
	// out of the session another one just created
	renewMutex sync.Mutex
}

func (p *Pihole) getSession() (string, time.Time) {
	p.sessionMutex.RLock()
	defer p.sessionMutex.RUnlock()

	return p.SID, p.ValidityTime
}

func (p *Pihole) setSession(sid string, validityTime time.Time) {
	p.sessionMutex.Lock()
	defer p.sessionMutex.Unlock()

	p.SID, p.ValidityTime = sid, validityTime
}

// renewSession logs out and logs in again. If staleSID isn't empty and the
// current SID is already a different one, another request renewed the session in
// the meantime and there is nothing to do.
func (p *Pihole) renewSession(staleSID string) error {
	p.renewMutex.Lock()
	defer p.renewMutex.Unlock()

	sid, _ := p.getSession()
	if staleSID != "" && sid != staleSID {
		return nil
	}

	if sid != "" {
		// a failing logout must not stop the login below, otherwise a session the
		// server no longer accepts would leave the client stuck until a restart
		if err := p.Logout(); err != nil {
			p.setSession("", time.Time{})
		}
	}

	return p.Login()
}

// newMutex serializes New() calls, otherwise concurrent requests can each
// build their own instance and race on the package singleton
var newMutex sync.Mutex

func New() (*Pihole, error) {
	newMutex.Lock()
	defer newMutex.Unlock()

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
