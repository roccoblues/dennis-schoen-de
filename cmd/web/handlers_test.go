package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttpsRedirect(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		tlsPort      string
		hostName     string
		wantStatus   int
		wantLocation string
	}{
		{
			name:       "no host",
			url:        "/",
			tlsPort:    "443",
			hostName:   "",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:         "no host; default hostname",
			url:          "/",
			tlsPort:      "443",
			hostName:     "www.test.com",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://www.test.com/",
		},
		{
			name:         "no host; default hostname; path",
			url:          "/foo/edit",
			tlsPort:      "443",
			hostName:     "www.test.com",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://www.test.com/foo/edit",
		},
		{
			name:         "non-standard port",
			url:          "/",
			tlsPort:      "4001",
			hostName:     "www.test.com",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://www.test.com:4001/",
		},
		{
			name:         "non-standard port; path",
			url:          "/foo/edit",
			tlsPort:      "4001",
			hostName:     "www.test.com",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://www.test.com:4001/foo/edit",
		},
		{
			name:         "with host",
			url:          "http://www.test.com",
			tlsPort:      "443",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://www.test.com",
		},
		{
			name:         "with host; non-standard port",
			url:          "http://www.test.com",
			tlsPort:      "4001",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://www.test.com:4001",
		},
		{
			name:         "with host; non-standard port; path",
			url:          "http://www.test.com/foo/edit",
			tlsPort:      "4001",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://www.test.com:4001/foo/edit",
		},
		{
			name:         "with host and hostname",
			url:          "http://www.test.com",
			tlsPort:      "443",
			hostName:     "www.something-else.com",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://www.something-else.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			r, err := http.NewRequest(http.MethodGet, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			app := newTestApplication(t)
			app.hostName = tt.hostName

			app.httpsRedirect(tt.tlsPort)(rr, r)

			rs := rr.Result()

			if rs.StatusCode != tt.wantStatus {
				t.Errorf("want %d; got %d", tt.wantStatus, rs.StatusCode)
			}

			location := rs.Header.Get("location")
			if location != tt.wantLocation {
				t.Errorf("want %q; got %q", tt.wantLocation, location)
			}
		})
	}
}

func TestHome(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if !strings.Contains(string(body), "my name is Dennis") {
		t.Errorf("want body to contain %q", "my name is Dennis")
	}
}

func TestResume(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/resume")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if !strings.Contains(string(body), "Awesome Job Title") {
		t.Errorf("want body to contain %q", "Awesome Job Title")
	}
}

func TestNotFound(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/foo/bar")

	if code != http.StatusNotFound {
		t.Errorf("want %d; got %d", http.StatusNotFound, code)
	}

	if !strings.Contains(string(body), "Not Found") {
		t.Errorf("want body to contain %q", "Not Found")
	}
}
