package linkwarden

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestGetLinksLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"response":[{"id":1},{"id":2},{"id":3}]}`)
	}))
	defer server.Close()

	l := &Linkwarden{InternalAddress: server.URL, Token: "token"}

	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{"limit greater than the number of links", 10, 3},
		{"limit smaller than the number of links", 2, 2},
		{"limit equal to the number of links", 3, 3},
		{"no limit", -1, 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			links, err := l.GetLinks(test.limit, "")
			if err != nil {
				t.Fatal(err)
			}
			if len(links) != test.want {
				t.Fatalf("expected %d links, got %d", test.want, len(links))
			}
		})
	}
}

func TestGetLinks(t *testing.T) {
	v, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get links", func(t *testing.T) {
		links, err := v.GetLinks(-1, "")
		if err != nil {
			t.Fatal(err)
		}

		for _, link := range links {
			if link.ID == 0 {
				t.Fatal("links with ID 0")
			}
		}
	})
}
