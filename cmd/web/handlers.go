package main

import (
	"net"
	"net/http"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{})
}

func (app *application) resume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.render(w, r, "resume.page.tmpl", &templateData{CV: app.cv})
}

func (app *application) httpsRedirect(tlsPort string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.URL
		u.Host = net.JoinHostPort(stripPort(r.Host), tlsPort)
		u.Scheme = "https"
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	}
}

// https://github.com/golang/go/commit/1ff19201fd898c3e1a0ed5d3458c81c1f062570b#diff-6c2d018290e298803c0c9419d8739885R971-R980
func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}
