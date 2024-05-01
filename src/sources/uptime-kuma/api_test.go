package uptimekuma

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

func TestUptimeKuma_GetStatusPageLastHeartbeats(t *testing.T) {
	u, err := New(config.GlobalConfigs.UptimeKumaConfigs.Address)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Test GetStatusPageLastHeartbeats", func(t *testing.T) {
		_, err := u.GetStatusPageLastUpDownCount("") // change to a valid slug
		if err != nil {
			t.Fatal(err)
		}
	})
}
