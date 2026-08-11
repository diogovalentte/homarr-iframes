package kavita

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// fakeKavita is a Kavita server that only accepts the token it handed out last,
// so an expired token has to be refreshed before the request goes through.
type fakeKavita struct {
	mutex        sync.Mutex
	rotations    int
	token        string
	refreshToken string
	logins       atomic.Int64
	refreshes    atomic.Int64
}

func (f *fakeKavita) rotate() (string, string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.rotations++
	f.token = fmt.Sprintf("token-%d", f.rotations)
	f.refreshToken = "refresh-" + f.token

	return f.token, f.refreshToken
}

func (f *fakeKavita) currentToken() string {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return f.token
}

func (f *fakeKavita) server() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/account/login", func(w http.ResponseWriter, _ *http.Request) {
		f.logins.Add(1)
		token, refreshToken := f.rotate()
		fmt.Fprintf(w, `{"token":%q,"refreshToken":%q}`, token, refreshToken)
	})

	mux.HandleFunc("/api/account/refresh-token", func(w http.ResponseWriter, _ *http.Request) {
		f.refreshes.Add(1)
		token, refreshToken := f.rotate()
		fmt.Fprintf(w, `{"token":%q,"refreshToken":%q}`, token, refreshToken)
	})

	mux.HandleFunc("/api/Server/media-errors", func(w http.ResponseWriter, r *http.Request) {
		if strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ") != f.currentToken() {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Fprint(w, `{"results":[{"comment":"broken file","createdUtc":"2026-08-10T00:00:00.0000000"}]}`)
	})

	return httptest.NewServer(mux)
}

// TestConcurrentRequestsRefreshOnce must be run with -race. Concurrent requests
// share one Kavita instance, so they read and write the tokens from different
// goroutines, and a single expired token must not cause one refresh per request.
func TestConcurrentRequestsRefreshOnce(t *testing.T) {
	fake := &fakeKavita{}
	server := fake.server()
	defer server.Close()

	k := &Kavita{InternalAddress: server.URL, Username: "user", Password: "pass"}
	if err := k.Login(); err != nil {
		t.Fatal(err)
	}

	// the server moves on to a new token, every in-flight request is now stale
	fake.rotate()

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := k.GetMediaErrors(); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	if refreshes := fake.refreshes.Load(); refreshes != 1 {
		t.Fatalf("expected the stale token to be refreshed once, got %d refreshes", refreshes)
	}
}

// TestRefreshWithoutTokensLogsIn covers the client recovering on its own when it
// has no tokens to refresh with, instead of erroring until it's restarted.
func TestRefreshWithoutTokensLogsIn(t *testing.T) {
	fake := &fakeKavita{}
	server := fake.server()
	defer server.Close()

	k := &Kavita{InternalAddress: server.URL, Username: "user", Password: "pass"}

	if err := k.RefreshCurrentToken(); err != nil {
		t.Fatal(err)
	}

	if logins := fake.logins.Load(); logins != 1 {
		t.Fatalf("expected 1 login, got %d", logins)
	}
	if token, refreshToken := k.getTokens(); token == "" || refreshToken == "" {
		t.Fatal("tokens are empty after refreshing")
	}
}
