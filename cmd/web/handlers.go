package main

import (
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
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			app.serverError(w, err)
			return
		}
		u := r.URL
		u.Host = net.JoinHostPort(host, tlsPort)
		u.Scheme = "https"
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	}
}
