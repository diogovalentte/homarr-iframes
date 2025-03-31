package backrest

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

func TestGetSummaryDashboard(t *testing.T) {
	b, err := New()
	if err != nil {
		t.Fatalf("error creating Backrest instance: %v", err)
	}
	_, err = b.GetSummaryDashboard()
	if err != nil {
		t.Fatalf("error getting summary dashboard: %v", err)
	}
}
