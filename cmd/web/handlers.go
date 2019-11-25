package main

import (
	"net/http")

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.tmpl", &templateData{})
}

func (app *application) resume(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "resume.page.tmpl", &templateData{})
}
