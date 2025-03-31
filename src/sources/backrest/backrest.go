package backrest

import (
	"fmt"
	"strings"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

var b *Backrest

type Backrest struct {
	Address         string
	InternalAddress string
	username        string
	password        string
}

func New() (*Backrest, error) {
	if b != nil {
		return b, nil
	}

	newR := &Backrest{}
	err := newR.Init()
	if err != nil {
		return nil, err
	}

	b = newR

	return b, nil
}

func (b *Backrest) Init() error {
	address, internalAddress, username, password := config.GlobalConfigs.Backrest.Address, config.GlobalConfigs.Backrest.InternalAddress, config.GlobalConfigs.Backrest.Username, config.GlobalConfigs.Backrest.Password
	if address == "" {
		return fmt.Errorf("BACKREST_ADDRESS variable should be set")
	}

	b.Address = strings.TrimSuffix(address, "/")
	if internalAddress == "" {
		b.InternalAddress = b.Address
	} else {
		b.InternalAddress = strings.TrimSuffix(internalAddress, "/")
	}
	b.username = username
	b.password = password

	return nil
}

func (b *Backrest) GetSummaryDashboard() (*SummaryDashboard, error) {
	var dashboard *SummaryDashboard
	err := b.baseRequest("POST", fmt.Sprintf("%s/v1.Backrest/GetSummaryDashboard", b.InternalAddress), nil, &dashboard)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}

type SummaryDashboard struct {
	PlanSummaries []PlanSummary `json:"planSummaries"`
}

type PlanSummary struct {
	ID                  string `json:"id"`
	BackupsFailed30days string `json:"backupsFailed30days"`
	RecentBackups       struct {
		FlowID      []string `json:"flowId"`
		TimeStampMs []string `json:"timeStampMs"`
		DurationMs  []string `json:"durationMs"`
		Status      []string `json:"status"`
	} `json:"recentBackups"`
}
