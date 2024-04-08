package config

import (
	"os"

	"github.com/joho/godotenv"
)

var GlobalConfigs Configs

type Configs struct {
	LinkwardenConfigs linkwardenConfigs
}

type linkwardenConfigs struct {
	Address string
	Token   string
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

	return nil
}
