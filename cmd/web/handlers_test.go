package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
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

			app := &application{
				errorLog: log.New(ioutil.Discard, "", 0),
				infoLog:  log.New(ioutil.Discard, "", 0),
				hostName: tt.hostName,
			}

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
