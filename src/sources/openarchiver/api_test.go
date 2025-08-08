package openarchiver

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
	o, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get links", func(t *testing.T) {
		ingestionSources, err := o.GetIngestionSources(-1)
		if err != nil {
			t.Fatal(err)
		}

		for _, ingestionSource := range ingestionSources {
			if ingestionSource.ID == "" {
				t.Fatal("ingestion source ID is empty")
			}
		}
	})
}
