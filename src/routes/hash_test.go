package routes_test

import (
	"net/http"
	"testing"
)

func TestGetHashes(t *testing.T) {
	t.Run("Get Linkwarden hash", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/hash/linkwarden", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Cinemark hash", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/hash/cinemark?theaterIds=2133", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Vikunja hash", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/hash/vikunja", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Media Releases hash", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/hash/media_releases", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Media Requests hash", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/hash/media_requests", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get UptimeKuma hash", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/hash/uptimekuma?slug=general", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
	t.Run("Get Alarms hash", func(t *testing.T) {
		r, err := requestHelper(http.MethodGet, "/v1/hash/alarms?alarms=sonarr,radarr,lidarr,prowlarr,kavita,pihole,speedtest-tracker,netdata,changedetectionio,kaizoku", nil)
		if err != nil {
			t.Fatal(err)
		}

		if r.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d: %s", r.Code, r.Body.String())
		}
	})
}
