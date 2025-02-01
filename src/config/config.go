package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	GlobalConfigs                            *Configs
	defaultChangeDetectionIOChangedLastHours = 24
)

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
	Jellyseerr              jellyseerrConfigs
	ChangeDetectionIO       changedetectionIOConfigs
}

type linkwardenConfigs struct {
	Address          string
	InternalAddress  string
	Token            string
	BackgroundImgURL string
}

type vikunjaConfigs struct {
	Address          string
	InternalAddress  string
	Token            string
	BackgroundImgURL string
}

type overseerrConfigs struct {
	Address         string
	InternalAddress string
	APIKey          string
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
	Token           string
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

type jellyseerrConfigs struct {
	Address         string
	InternalAddress string
	APIKey          string
}

type changedetectionIOConfigs struct {
	Address          string
	InternalAddress  string
	APIKey           string
	ChangedLastHours int
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
	GlobalConfigs.Linkwarden.BackgroundImgURL = os.Getenv("LINKWARDEN_BACKGROUND_IMG_URL")

	GlobalConfigs.Vikunja.Address = os.Getenv("VIKUNJA_ADDRESS")
	GlobalConfigs.Vikunja.InternalAddress = os.Getenv("INTERNAL_VIKUNJA_ADDRESS")
	GlobalConfigs.Vikunja.Token = os.Getenv("VIKUNJA_TOKEN")
	GlobalConfigs.Vikunja.BackgroundImgURL = os.Getenv("VIKUNJA_BACKGROUND_IMG_URL")

	GlobalConfigs.Overseerr.Address = os.Getenv("OVERSEERR_ADDRESS")
	GlobalConfigs.Overseerr.InternalAddress = os.Getenv("INTERNAL_OVERSEERR_ADDRESS")
	GlobalConfigs.Overseerr.APIKey = os.Getenv("OVERSEERR_API_KEY")

	GlobalConfigs.Jellyseerr.Address = os.Getenv("JELLYSEERR_ADDRESS")
	GlobalConfigs.Jellyseerr.InternalAddress = os.Getenv("INTERNAL_JELLYSEERR_ADDRESS")
	GlobalConfigs.Jellyseerr.APIKey = os.Getenv("JELLYSEERR_API_KEY")

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
	GlobalConfigs.SpeedTestTrackerConfigs.Token = os.Getenv("SPEEDTEST_TRACKER_TOKEN")

	GlobalConfigs.Pihole.Address = os.Getenv("PIHOLE_ADDRESS")
	GlobalConfigs.Pihole.InternalAddress = os.Getenv("INTERNAL_PIHOLE_ADDRESS")
	GlobalConfigs.Pihole.Token = os.Getenv("PIHOLE_TOKEN")

	GlobalConfigs.Kavita.Address = os.Getenv("KAVITA_ADDRESS")
	GlobalConfigs.Kavita.InternalAddress = os.Getenv("INTERNAL_KAVITA_ADDRESS")
	GlobalConfigs.Kavita.Username = os.Getenv("KAVITA_USERNAME")
	GlobalConfigs.Kavita.Password = os.Getenv("KAVITA_PASSWORD")

	GlobalConfigs.Kaizoku.Address = os.Getenv("KAIZOKU_ADDRESS")
	GlobalConfigs.Kaizoku.InternalAddress = os.Getenv("INTERNAL_KAIZOKU_ADDRESS")

	GlobalConfigs.ChangeDetectionIO.Address = os.Getenv("CHANGEDETECTIONIO_ADDRESS")
	GlobalConfigs.ChangeDetectionIO.InternalAddress = os.Getenv("INTERNAL_CHANGEDETECTIONIO_ADDRESS")
	GlobalConfigs.ChangeDetectionIO.APIKey = os.Getenv("CHANGEDETECTIONIO_API_KEY")
	changedLastHours := os.Getenv("CHANGEDETECTIONIO_CHANGED_LAST_HOURS")
	if changedLastHours == "" {
		GlobalConfigs.ChangeDetectionIO.ChangedLastHours = defaultChangeDetectionIOChangedLastHours
	} else {
		GlobalConfigs.ChangeDetectionIO.ChangedLastHours, err = strconv.Atoi(changedLastHours)
		if err != nil {
			return err
		}
	}

	return nil
}
