package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/markbates/pkger"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/resume", http.HandlerFunc(app.resume))

	fileServer := http.FileServer(pkger.Dir("/ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
