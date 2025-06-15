package routes_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	api "github.com/diogovalentte/homarr-iframes/src"
	"github.com/diogovalentte/homarr-iframes/src/config"
)

func setup() error {
	envFilePath := "../../.env.test"
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

func TestGetIFrames(t *testing.T) {
	t.Run("Get Linkwarden iFrame", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/iframe/linkwarden", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Cinemark iFrame", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/iframe/cinemark?theaterIds=2133", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Vikunja iFrame", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/iframe/vikunja", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Overseerr iFrame", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/iframe/overseerr", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusMovedPermanently {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Media Releases iFrame", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/iframe/media_releases", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Media Requests iFrame", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/iframe/media_requests", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get UptimeKuma iFrame", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/iframe/uptimekuma?slug=general", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Alarms iFrame", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/iframe/alarms?alarms=sonarr,radarr,lidarr,prowlarr,kavita,pihole,speedtest-tracker,netdata,changedetectionio,kaizoku", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Netdata iFrame", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/iframe/netdata", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusMovedPermanently {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
}

func TestSetVikunjaTaskDone(t *testing.T) {
	t.Run("Set Vikunja Task Done", func(t *testing.T) {
		r, err := requestHelper(http.MethodPatch, "/v1/iframe/vikunja/set_task_done?taskId=1", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}

		rb := map[string]string{}
		err = json.Unmarshal(r.Body.Bytes(), &rb)
		if err != nil {
			t.Fatal(err)
		}
		if rb["message"] != "Task done" {
			t.Fatalf("expected message 'Task done', got %s", rb["message"])
		}
	})
}

func requestHelper(method, url string, target any) (*httptest.ResponseRecorder, error) {
	r := httptest.NewRecorder()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	router := api.SetupRouter()
	router.ServeHTTP(r, req)

	return r, nil
}
