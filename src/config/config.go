package config

import (
	"os"

	"github.com/joho/godotenv"
)

var GlobalConfigs *Configs

type Configs struct {
	Linkwarden        linkwardenConfigs
	Vikunja           vikunjaConfigs
	Overseerr         overseerrConfigs
	Sonarr            sonarrConfigs
	Radarr            radarrConfigs
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

type overseerrConfigs struct {
	Address string
	Token   string
}

type sonarrConfigs struct {
	Address string
	APIKey  string
}

type radarrConfigs struct {
	Address string
	APIKey  string
}

type uptimeKumaConfigs struct {
	Address string
}

func SetConfigs(filePath string) error {
	GlobalConfigs = &Configs{}

	var err error
	if filePath != "" {
		err = godotenv.Load(filePath)
		if err != nil {
			return err
		}
	}

	GlobalConfigs.Linkwarden.Address = os.Getenv("LINKWARDEN_ADDRESS")
	GlobalConfigs.Linkwarden.Token = os.Getenv("LINKWARDEN_TOKEN")

	GlobalConfigs.Vikunja.Address = os.Getenv("VIKUNJA_ADDRESS")
	GlobalConfigs.Vikunja.Token = os.Getenv("VIKUNJA_TOKEN")

	GlobalConfigs.Overseerr.Address = os.Getenv("OVERSEERR_ADDRESS")
	GlobalConfigs.Overseerr.Token = os.Getenv("OVERSEERR_TOKEN")

	GlobalConfigs.Sonarr.Address = os.Getenv("SONARR_ADDRESS")
	GlobalConfigs.Sonarr.APIKey = os.Getenv("SONARR_API_KEY")

	GlobalConfigs.Radarr.Address = os.Getenv("RADARR_ADDRESS")
	GlobalConfigs.Radarr.APIKey = os.Getenv("RADARR_API_KEY")

	GlobalConfigs.UptimeKumaConfigs.Address = os.Getenv("UPTIMEKUMA_ADDRESS")

	return nil
}
