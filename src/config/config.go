package config

import (
	"os"

	"github.com/joho/godotenv"
)

var GlobalConfigs Configs

type Configs struct {
	LinkwardenConfigs linkwardenConfigs
	VikunjaConfigs    vikunjaConfigs
}

type linkwardenConfigs struct {
	Address string
	Token   string
}

type vikunjaConfigs struct {
	Address  string
	Username string
	Password string
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
	GlobalConfigs.VikunjaConfigs.Username = os.Getenv("VIKUNJA_USERNAME")
	GlobalConfigs.VikunjaConfigs.Password = os.Getenv("VIKUNJA_PASSWORD")

	return nil
}
