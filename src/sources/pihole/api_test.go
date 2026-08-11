package pihole

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// fakePihole is a Pi-hole v6 server that only accepts the session it handed out
// last, so an expired SID has to be renewed before the request goes through.
type fakePihole struct {
	mutex     sync.Mutex
	sessions  int
	sid       string
	logins    atomic.Int64
	logouts   atomic.Int64
	messages  atomic.Int64
	lastBody  string
	bodyMutex sync.Mutex
}

func (f *fakePihole) rotate() string {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.sessions++
	f.sid = fmt.Sprintf("sid-%d", f.sessions)

	return f.sid
}

func (f *fakePihole) currentSID() string {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return f.sid
}

func (f *fakePihole) server() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth", func(w http.ResponseWriter, _ *http.Request) {
		f.logins.Add(1)
		fmt.Fprintf(w, `{"session":{"valid":true,"sid":%q,"validity":1800}}`, f.rotate())
	})

	mux.HandleFunc("DELETE /api/auth", func(w http.ResponseWriter, _ *http.Request) {
		f.logouts.Add(1)
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("/api/info/messages", func(w http.ResponseWriter, r *http.Request) {
		f.messages.Add(1)
		if r.Header.Get("X-FTL-SID") != f.currentSID() {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"error":{"message":"Unauthorized"}}`)
			return
		}
		body, _ := io.ReadAll(r.Body)
		f.bodyMutex.Lock()
		f.lastBody = string(body)
		f.bodyMutex.Unlock()
		fmt.Fprint(w, `{"messages":[{"type":"test","plain":"a message"}]}`)
	})

	return httptest.NewServer(mux)
}

func newFake(t *testing.T) (*Pihole, *fakePihole) {
	t.Helper()

	fake := &fakePihole{}
	server := fake.server()
	t.Cleanup(server.Close)

	p := &Pihole{InternalAddress: server.URL, Password: "password"}
	if err := p.Login(); err != nil {
		t.Fatal(err)
	}

	return p, fake
}

// TestConcurrentRequestsRenewOnce must be run with -race. Concurrent requests
// share one Pi-hole instance, so they read and write the session from different
// goroutines, and an expired session must not be renewed once per request.
func TestConcurrentRequestsRenewOnce(t *testing.T) {
	p, fake := newFake(t)

	// the server moves on to a new session, every in-flight request is now stale
	fake.rotate()

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := p.GetMessages(); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	if logins := fake.logins.Load(); logins != 2 {
		t.Fatalf("expected the initial login plus a single renewal, got %d logins", logins)
	}
}

// TestRetryReplaysTheBody covers the retry after an unauthorized response, which
// used to reuse the already consumed body reader and send an empty body.
func TestRetryReplaysTheBody(t *testing.T) {
	p, fake := newFake(t)

	// force the retry path: the session the request uses is no longer accepted
	fake.rotate()

	var target map[string]any
	body := `{"field":"value"}`
	if err := p.baseRequest(http.MethodPost, p.InternalAddress+"/api/info/messages", strings.NewReader(body), &target, 1); err != nil {
		t.Fatal(err)
	}

	fake.bodyMutex.Lock()
	defer fake.bodyMutex.Unlock()
	if fake.lastBody != body {
		t.Fatalf("expected the retry to send %q, got %q", body, fake.lastBody)
	}
}

// TestEmptySessionRenewsBeforeRequesting covers an empty session being renewed
// upfront. It used to skip the renewal, send a request with an empty session
// header and only recover from the unauthorized response it got back, wasting a
// round trip on every request until the session was renewed for another reason.
func TestEmptySessionRenewsBeforeRequesting(t *testing.T) {
	p, fake := newFake(t)

	_, validityTime := p.getSession()
	p.setSession("", validityTime)

	if _, err := p.GetMessages(); err != nil {
		t.Fatal(err)
	}
	if logins := fake.logins.Load(); logins != 2 {
		t.Fatalf("expected the initial login plus a login for the empty session, got %d logins", logins)
	}
	if messages := fake.messages.Load(); messages != 1 {
		t.Fatalf("expected the session to be renewed before requesting, got %d requests", messages)
	}
}
