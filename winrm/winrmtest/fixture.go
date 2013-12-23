package winrmtest

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

type Fixture struct {
	mux      *http.ServeMux
	server   *httptest.Server
	Endpoint string
}

func NewFixture() *Fixture {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	return &Fixture{mux, server, server.URL}
}

func (f *Fixture) HandleFunc(handler func(w http.ResponseWriter, r *Request)) {
	f.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if canHandle(r) {
			handler(w, newRequest(r))
		}

		go f.server.Close()
	})
}

func canHandle(r *http.Request) bool {
	if r.Method != "POST" {
		return false
	}

	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/soap+xml") {
		return false
	}

	return true
}
