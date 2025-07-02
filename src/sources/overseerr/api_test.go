package overseerr

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
	o, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get requests", func(t *testing.T) {
		requests, err := o.GetRequests(-1, "", "", 0)
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
	o, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get movie", func(t *testing.T) {
		media, err := o.GetMedia(-1, "", "") // civil war
		if err != nil {
			t.Fatal(err)
		}

		if len(media) == 0 {
			t.Fatal("empty media list")
		}
	})
}

func TestGetMovieTV(t *testing.T) {
	t.Run("get movie", func(t *testing.T) {
		media, err := o.GetMovie(929590) // civil war
		if err != nil {
			t.Fatal(err)
		}

		if media.ID == 0 {
			t.Fatal("media with ID 0")
		}
	})

	t.Run("get tv show", func(t *testing.T) {
		media, err := o.GetTV(2316) // the office
		if err != nil {
			t.Fatal(err)
		}

		if media.ID == 0 {
			t.Fatal("media with ID 0")
		}
	})
}

func TestGetIframeData(t *testing.T) {
	o, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get iframe requests data", func(t *testing.T) {
		iframeData, err := o.GetIframeData(-1, "", "", 0, false)
		if err != nil {
			t.Fatal(err)
		}

		if len(iframeData) == 0 {
			t.Fatal("empty iframe requests data")
		}
	})

	t.Run("get iframe media data", func(t *testing.T) {
		iframeData, err := o.GetIframeData(-1, "", "", 0, true)
		if err != nil {
			t.Fatal(err)
		}

		if len(iframeData) == 0 {
			t.Fatal("empty iframe media data")
		}
	})
}
