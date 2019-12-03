package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectHostName(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		hostName     string
		wantStatus   int
		wantLocation string
	}{
		{
			name:       "no host; no hostname",
			url:        "/",
			hostName:   "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "host; no hostname",
			url:        "http://www.test.com",
			hostName:   "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "host; no hostname; non-standard port",
			url:        "http://www.test.com:4000",
			hostName:   "",
			wantStatus: http.StatusOK,
		},
		{
			name:         "no host; hostname",
			url:          "/",
			hostName:     "www.test.com",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "http://www.test.com/",
		},
		{
			name:       "same host and hostname",
			url:        "http://www.test.com",
			hostName:   "www.test.com",
			wantStatus: http.StatusOK,
		},
		{
			name:       "same host and hostname; non-standard port",
			url:        "http://www.test.com:4000",
			hostName:   "www.test.com",
			wantStatus: http.StatusOK,
		},
		{
			name:       "same host and hostname; with path",
			url:        "http://www.test.com/foo/edit",
			hostName:   "www.test.com",
			wantStatus: http.StatusOK,
		},
		{
			name:       "same host and hostname; with path; https",
			url:        "https://www.test.com/foo/edit",
			hostName:   "www.test.com",
			wantStatus: http.StatusOK,
		},
		{
			name:         "different host and hostname",
			url:          "http://www.foo.com",
			hostName:     "www.test.com",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "http://www.test.com",
		},
		{
			name:         "different host and hostname; non-standard port",
			url:          "http://www.foo.com:4000",
			hostName:     "www.test.com",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "http://www.test.com:4000",
		},
		{
			name:         "different host and hostname; non-standard port; https",
			url:          "https://www.foo.com:4000",
			hostName:     "www.test.com",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://www.test.com:4000",
		},
		{
			name:         "different host and hostname; with path",
			url:          "http://www.foo.com/foo/edit",
			hostName:     "www.test.com",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "http://www.test.com/foo/edit",
		},
		{
			name:         "different host and hostname; with path; https",
			url:          "https://www.foo.com/foo/edit",
			hostName:     "www.test.com",
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://www.test.com/foo/edit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)

			app := &application{
				errorLog: log.New(ioutil.Discard, "", 0),
				infoLog:  log.New(ioutil.Discard, "", 0),
				hostName: tt.hostName,
			}

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})

			app.redirectHostName(next).ServeHTTP(rr, r)

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

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}

	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", xssProtection)
	}

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
