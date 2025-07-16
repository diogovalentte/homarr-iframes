package jellyfin

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

func TestGetLatestItems(t *testing.T) {
	j, err := New()
	if err != nil {
		t.Fatalf("error creating Jellyfin instance: %v", err)
	}

	items, err := j.GetLatestItems(20, 100, j.userId, "", "Movie,Episode")
	if err != nil {
		t.Fatalf("error getting latest items: %v", err)
	}

	t.Logf("Retrieved %d items from Jellyfin", len(items))
	if len(items) > 0 {
		t.Logf("First item: %+v", items[0])
	}
}
