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

func TestGetSessions(t *testing.T) {
	j, err := New()
	if err != nil {
		t.Fatalf("error creating Jellyfin instance: %v", err)
	}

	sessions, err := j.GetSessions(20, 3600)
	if err != nil {
		t.Fatalf("error getting sessions: %v", err)
	}

	t.Logf("Retrieved %d active sessions from Jellyfin", len(sessions))
	if len(sessions) > 0 {
		t.Logf("First session: %+v", sessions[0])
		if sessions[0].NowPlayingItem != nil {
			t.Logf("Now playing item: %+v", sessions[0].NowPlayingItem)
		}
	}
}
