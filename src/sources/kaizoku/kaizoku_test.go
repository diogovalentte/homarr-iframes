package kaizoku

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

func TestGetQueues(t *testing.T) {
	k, err := New()
	if err != nil {
		t.Fatalf("error creating Kaizoku instance: %v", err)
	}
	_, err = k.GetQueues()
	if err != nil {
		t.Fatalf("error getting queues: %v", err)
	}
}
