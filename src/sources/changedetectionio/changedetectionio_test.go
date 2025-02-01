package changedetectionio

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

func TestGetWatches(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("error creating changedetectionio instance: %v", err)
	}
	_, err = c.GetWatches()
	if err != nil {
		t.Fatalf("error getting watches: %v", err)
	}
}
