package jellyseerr

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

func TestGetRequests(t *testing.T) {
	j, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get requests", func(t *testing.T) {
		requests, err := j.GetRequests(-1, "", "", 0)
		if err != nil {
			t.Fatal(err)
		}

		for _, request := range requests {
			if request.ID == 0 {
				t.Fatal("task with ID 0")
			}
		}
	})
}

func TestGetMedia(t *testing.T) {
	j, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get movie", func(t *testing.T) {
		media, err := j.GetMedia("movie", 929590) // civil war
		if err != nil {
			t.Fatal(err)
		}

		if media.ID == 0 {
			t.Fatal("media with ID 0")
		}
	})

	t.Run("get tv show", func(t *testing.T) {
		media, err := j.GetMedia("tv", 2316) // the office
		if err != nil {
			t.Fatal(err)
		}

		if media.ID == 0 {
			t.Fatal("media with ID 0")
		}
	})
}

func TestGetIframeData(t *testing.T) {
	j, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get iframe data", func(t *testing.T) {
		iframeData, err := j.GetIframeData(-1, "", "", 0)
		if err != nil {
			t.Fatal(err)
		}

		if len(iframeData) == 0 {
			t.Fatal("empty iframe data")
		}
	})
}
