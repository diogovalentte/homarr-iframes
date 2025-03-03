package pihole

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

func TestGetMessages(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("error creating Pi-hole instance: %v", err)
	}
	_, err = p.GetMessages()
	if err != nil {
		t.Fatalf("error getting messages: %v", err)
	}
}

func TestLogin(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatalf("error creating Pi-hole instance: %v", err)
	}
}

func TestLogout(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("error creating Pi-hole instance: %v", err)
	}
	err = p.Logout()
	if err != nil {
		t.Fatalf("error logging out: %v", err)
	}
}
