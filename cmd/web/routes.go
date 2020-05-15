package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.defaultMiddleware(app.home))
	mux.HandleFunc("/resume", app.defaultMiddleware(app.resume))

	// and some fun: https://twitter.com/LiamHammett/status/1260984553570570240
	mux.HandleFunc("/.env", app.defaultMiddleware(app.redirectFun))
	mux.HandleFunc("/wp-login.php", app.defaultMiddleware(app.redirectFun))
	mux.HandleFunc("/wp-admin", app.defaultMiddleware(app.redirectFun))
	mux.HandleFunc("/xmlrpc.php", app.defaultMiddleware(app.redirectFun))
	mux.HandleFunc("/cgi-bin/mainfunction.cgi", app.defaultMiddleware(app.redirectFun))
	mux.HandleFunc("/owa/auth/logon.aspx", app.defaultMiddleware(app.redirectFun))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

func (app *application) defaultMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return app.recoverPanic(app.logRequest(secureHeaders(app.redirectHostName(next))))
}
