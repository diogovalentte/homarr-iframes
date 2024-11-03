package kaizoku

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var k *Kaizoku

type Kaizoku struct {
	Address         string
	InternalAddress string
}

func New() (*Kaizoku, error) {
	if k != nil {
		return k, nil
	}

	newR := &Kaizoku{}
	err := newR.Init()
	if err != nil {
		return nil, err
	}

	k = newR

	return k, nil
}

func (k *Kaizoku) Init() error {
	address, internalAddress := config.GlobalConfigs.Kaizoku.Address, config.GlobalConfigs.Kaizoku.InternalAddress
	if address == "" {
		return fmt.Errorf("KAIZOKU_ADDRESS variable should be set")
	}

	k.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		k.InternalAddress = k.Address
	} else {
		k.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}

	return nil
}

func (k *Kaizoku) GetQueues() ([]*Queue, error) {
	url := fmt.Sprintf("%s/bull/queues/api/queues", k.Address)
	var queues getQueuesResponse
	err := baseRequest(http.MethodGet, url, nil, &queues)
	if err != nil {
		return nil, err
	}

	return queues.Queues, nil
}

type getQueuesResponse struct {
	Queues []*Queue `json:"queues"`
}

type Queue struct {
	Name   string `json:"name"`
	Counts struct {
		Active          int `json:"active"`
		Completed       int `json:"completed"`
		Delayed         int `json:"delayed"`
		Failed          int `json:"failed"`
		Paused          int `json:"paused"`
		Waiting         int `json:"waiting"`
		WaitingChildren int `json:"waiting-children"`
	} `json:"counts"`
}
