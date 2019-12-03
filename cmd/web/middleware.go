package main

import (
	"fmt"
	"net"

	"net/http"
)

func (app *application) recoverPanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next(w, r)
	}
}

func secureHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000")

		next(w, r)
	}
}

func (app *application) logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next(w, r)
	}
}

func (app *application) redirectHostName(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if app.hostName == "" || stripPort(r.Host) == app.hostName {
			next(w, r)
			return
		}

		port := portOnly(r.Host)

		u := r.URL
		if r.TLS != nil {
			u.Scheme = "https"
			if port == "" {
				port = "443"
			}
		} else {
			u.Scheme = "http"
			if port == "" {
				port = "80"
			}
		}
		if port == "80" || port == "443" {
			u.Host = app.hostName
		} else {
			u.Host = net.JoinHostPort(app.hostName, port)
		}

		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	}
}
