package config

import (
	"os"

	"github.com/joho/godotenv"
)

var GlobalConfigs Configs

type Configs struct {
	LinkwardenConfigs linkwardenConfigs
	VikunjaConfigs    vikunjaConfigs
	UptimeKumaConfigs uptimeKumaConfigs
}

type linkwardenConfigs struct {
	Address string
	Token   string
}

type vikunjaConfigs struct {
	Address string
	Token   string
}

type uptimeKumaConfigs struct {
	Address string
}

func SetConfigs(filePath string) error {
	if filePath != "" {
		err := godotenv.Load(filePath)
		if err != nil {
			return err
		}
	}

	GlobalConfigs.LinkwardenConfigs.Address = os.Getenv("LINKWARDEN_ADDRESS")
	GlobalConfigs.LinkwardenConfigs.Token = os.Getenv("LINKWARDEN_TOKEN")

	GlobalConfigs.VikunjaConfigs.Address = os.Getenv("VIKUNJA_ADDRESS")
	GlobalConfigs.VikunjaConfigs.Token = os.Getenv("VIKUNJA_TOKEN")

	GlobalConfigs.UptimeKumaConfigs.Address = os.Getenv("UPTIMEKUMA_ADDRESS")

	return nil
}
