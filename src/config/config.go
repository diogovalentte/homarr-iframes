package config

import (
	"os"

	"github.com/joho/godotenv"
)

var GlobalConfigs *Configs

type Configs struct {
	Linkwarden              linkwardenConfigs
	Vikunja                 vikunjaConfigs
	Overseerr               overseerrConfigs
	Sonarr                  sonarrConfigs
	Radarr                  radarrConfigs
	Prowlarr                prowlarrConfigs
	UptimeKumaConfigs       uptimeKumaConfigs
	NetdataConfigs          netdataConfigs
	SpeedTestTrackerConfigs speedTestTrackerConfigs
	Pihole                  piholeConfigs
	Kavita                  kavitaConfigs
	Kaizoku                 kaizokuConfigs
}

type linkwardenConfigs struct {
	Address         string
	InternalAddress string
	Token           string
}

type vikunjaConfigs struct {
	Address         string
	InternalAddress string
	Token           string
}

type overseerrConfigs struct {
	Address         string
	InternalAddress string
	Token           string
}

type sonarrConfigs struct {
	Address         string
	InternalAddress string
	APIKey          string
}

type radarrConfigs struct {
	Address         string
	InternalAddress string
	APIKey          string
}

type prowlarrConfigs struct {
	Address         string
	InternalAddress string
	APIKey          string
}

type uptimeKumaConfigs struct {
	Address         string
	InternalAddress string
}

type netdataConfigs struct {
	Address         string
	InternalAddress string
	Token           string
}

type speedTestTrackerConfigs struct {
	Address         string
	InternalAddress string
}

type piholeConfigs struct {
	Address         string
	InternalAddress string
	Token           string
}

type kavitaConfigs struct {
	Address         string
	InternalAddress string
	Username        string
	Password        string
}

type kaizokuConfigs struct {
	Address         string
	InternalAddress string
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
	GlobalConfigs.Linkwarden.InternalAddress = os.Getenv("INTERNAL_LINKWARDEN_ADDRESS")
	GlobalConfigs.Linkwarden.Token = os.Getenv("LINKWARDEN_TOKEN")

	GlobalConfigs.Vikunja.Address = os.Getenv("VIKUNJA_ADDRESS")
	GlobalConfigs.Vikunja.InternalAddress = os.Getenv("INTERNAL_VIKUNJA_ADDRESS")
	GlobalConfigs.Vikunja.Token = os.Getenv("VIKUNJA_TOKEN")

	GlobalConfigs.Overseerr.Address = os.Getenv("OVERSEERR_ADDRESS")
	GlobalConfigs.Overseerr.InternalAddress = os.Getenv("INTERNAL_OVERSEERR_ADDRESS")
	GlobalConfigs.Overseerr.Token = os.Getenv("OVERSEERR_TOKEN")

	GlobalConfigs.Sonarr.Address = os.Getenv("SONARR_ADDRESS")
	GlobalConfigs.Sonarr.InternalAddress = os.Getenv("INTERNAL_SONARR_ADDRESS")
	GlobalConfigs.Sonarr.APIKey = os.Getenv("SONARR_API_KEY")

	GlobalConfigs.Radarr.Address = os.Getenv("RADARR_ADDRESS")
	GlobalConfigs.Radarr.InternalAddress = os.Getenv("INTERNAL_RADARR_ADDRESS")
	GlobalConfigs.Radarr.APIKey = os.Getenv("RADARR_API_KEY")

	GlobalConfigs.Prowlarr.Address = os.Getenv("PROWLARR_ADDRESS")
	GlobalConfigs.Prowlarr.InternalAddress = os.Getenv("INTERNAL_PROWLARR_ADDRESS")
	GlobalConfigs.Prowlarr.APIKey = os.Getenv("PROWLARR_API_KEY")

	GlobalConfigs.UptimeKumaConfigs.Address = os.Getenv("UPTIMEKUMA_ADDRESS")
	GlobalConfigs.UptimeKumaConfigs.InternalAddress = os.Getenv("INTERNAL_UPTIMEKUMA_ADDRESS")

	GlobalConfigs.NetdataConfigs.Address = os.Getenv("NETDATA_ADDRESS")
	GlobalConfigs.NetdataConfigs.InternalAddress = os.Getenv("INTERNAL_NETDATA_ADDRESS")
	GlobalConfigs.NetdataConfigs.Token = os.Getenv("NETDATA_TOKEN")

	GlobalConfigs.SpeedTestTrackerConfigs.Address = os.Getenv("SPEEDTEST_TRACKER_ADDRESS")
	GlobalConfigs.SpeedTestTrackerConfigs.InternalAddress = os.Getenv("INTERNAL_SPEEDTEST_TRACKER_ADDRESS")

	GlobalConfigs.Pihole.Address = os.Getenv("PIHOLE_ADDRESS")
	GlobalConfigs.Pihole.InternalAddress = os.Getenv("INTERNAL_PIHOLE_ADDRESS")
	GlobalConfigs.Pihole.Token = os.Getenv("PIHOLE_TOKEN")

	GlobalConfigs.Kavita.Address = os.Getenv("KAVITA_ADDRESS")
	GlobalConfigs.Kavita.InternalAddress = os.Getenv("INTERNAL_KAVITA_ADDRESS")
	GlobalConfigs.Kavita.Username = os.Getenv("KAVITA_USERNAME")
	GlobalConfigs.Kavita.Password = os.Getenv("KAVITA_PASSWORD")

	GlobalConfigs.Kaizoku.Address = os.Getenv("KAIZOKU_ADDRESS")
	GlobalConfigs.Kaizoku.InternalAddress = os.Getenv("INTERNAL_KAIZOKU_ADDRESS")

	return nil
}
