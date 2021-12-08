package registry

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type Server struct {
	mu     *sync.RWMutex
	called bool
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.called = true
	w.Write([]byte("ok"))
}

func TestMain(t *testing.T) {
	server := &Server{
		mu:     &sync.RWMutex{},
		called: false,
	}
	srv := httptest.NewServer(server)
	defer srv.Close()

	r := New(srv.URL, "test", "host", "80")
	r.SetDelay(1)

	go r.Start()
	time.Sleep(time.Millisecond * 200)
	r.Stop()

	server.mu.RLock()
	defer server.mu.RUnlock()

	if server.called != true {
		t.Errorf("expected true got false")
	}
}
