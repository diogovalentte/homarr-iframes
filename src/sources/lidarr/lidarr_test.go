package lidarr

import (
	"fmt"
	"os"
	"testing"
	"time"

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

func TestGetCalendar(t *testing.T) {
	l, err := New()
	if err != nil {
		t.Fatalf("error creating Lidarr instance: %v", err)
	}
	_, err = l.GetCalendar(false, time.Now(), time.Now().AddDate(0, 0, 1))
	if err != nil {
		t.Fatalf("error getting calendar: %v", err)
	}
}

func TestGetHealth(t *testing.T) {
	l, err := New()
	if err != nil {
		t.Fatalf("error creating Lidarr instance: %v", err)
	}
	_, err = l.GetHealth()
	if err != nil {
		t.Fatalf("error getting health: %v", err)
	}
}
