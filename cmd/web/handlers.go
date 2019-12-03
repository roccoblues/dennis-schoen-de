package main

import (
	"fmt"
	"net"
	"net/http"
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
		host := app.hostName // default to the configured hostname
		if host == "" {      // try to get the hostname from request
			host = stripPort(r.Host)
		}
		if host == "" {
			// without a hostname we can't build a redirect
			app.serverError(w, fmt.Errorf("request host missing"))
			return
		}

		u := r.URL
		if tlsPort != "443" {
			u.Host = net.JoinHostPort(host, tlsPort)
		} else {
			u.Host = host
		}
		u.Scheme = "https"

		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	}
}
