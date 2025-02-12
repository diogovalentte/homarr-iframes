package alarms

import (
	"fmt"
	"os"
	"testing"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

func setup() error {
	envFilePath := "../../../.env.test"
	err := config.SetConfigs(envFilePath)
	if err != nil {
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetAlarms(t *testing.T) {
	alarms := Alarms{}
	t.Run("get alarms", func(t *testing.T) {
		_, err := alarms.GetAlarms([]string{"netdata", "prowlarr", "radarr", "sonarr", "speedtest-tracker"}, -1, false, true)
		if err != nil {
			t.Fatal(err)
		}
	})
}
