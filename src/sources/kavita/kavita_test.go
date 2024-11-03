package kavita

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

func TestLogin(t *testing.T) {
	k, err := New()
	if err != nil {
		t.Fatalf("error creating Kavita instance: %v", err)
	}
	err = k.Login()
	if err != nil {
		t.Fatalf("error logging in: %v", err)
	}
}

func TestRefreshToken(t *testing.T) {
	k, err := New()
	if err != nil {
		t.Fatalf("error creating Kavita instance: %v", err)
	}
	err = k.RefreshCurrentToken()
	if err != nil {
		t.Fatalf("error refreshing token: %v", err)
	}
}

func TestGetMediaErrors(t *testing.T) {
	k, err := New()
	if err != nil {
		t.Fatalf("error creating Kavita instance: %v", err)
	}
	_, err = k.GetMediaErrors()
	if err != nil {
		t.Fatalf("error getting media errors: %v", err)
	}
}
