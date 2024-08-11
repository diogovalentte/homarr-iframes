package linkwarden

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

func TestGetLinks(t *testing.T) {
	v, err := New(config.GlobalConfigs.Linkwarden.Address, config.GlobalConfigs.Linkwarden.InternalAddress, config.GlobalConfigs.Linkwarden.Token)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get links", func(t *testing.T) {
		links, err := v.GetLinks(-1, "")
		if err != nil {
			t.Fatal(err)
		}

		for _, link := range links {
			if link.ID == 0 {
				t.Fatal("links with ID 0")
			}
		}
	})
}
