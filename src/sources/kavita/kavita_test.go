package kavita

import (
	"encoding/json"
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

func TestGetMediaErrorsUnmarshal(t *testing.T) {
	// Kavita <= 0.9.0 returns a paginated object, 0.9.1 returns a plain array
	bodies := []string{
		`{"results":[{"comment":"This archive cannot be read or not supported","createdUtc":"2026-09-01T10:55:57.9649185"}]}`,
		`[{"extension":"CBZ","filePath":"/a.cbz","comment":"This archive cannot be read or not supported","details":"","createdUtc":"2026-09-01T10:55:57.9649185"}]`,
	}

	for _, body := range bodies {
		var errors MediaErrorResults
		if err := json.Unmarshal([]byte(body), &errors); err != nil {
			t.Fatalf("error unmarshaling %s: %s", body, err)
		}
		if len(errors.Results) != 1 {
			t.Fatalf("expected 1 media error from %s, got %d", body, len(errors.Results))
		}
		if errors.Results[0].Comment != "This archive cannot be read or not supported" {
			t.Fatalf("unexpected comment from %s: %s", body, errors.Results[0].Comment)
		}
	}
}
