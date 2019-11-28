package main

import (
	"net/http"

	"github.com/markbates/pkger"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/resume", app.resume)

	fileServer := http.FileServer(pkger.Dir("/ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
