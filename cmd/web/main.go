package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/roccoblues/dennis-schoen.de/pkg/models"
	"github.com/roccoblues/dennis-schoen.de/pkg/models/yml"
)

var version string

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	cv            *models.CV
	templateCache map[string]*template.Template
}

func main() {
	fs := flag.NewFlagSet("web", flag.ExitOnError)
	var (
		sslCert     = flag.String("ssl-cert", "./tls/localhost.pem", "SSL Certificate")
		sslKey      = flag.String("ssl-key", "./tls/localhost-key.pem", "SSL Key")
		httpsAddr   = fs.String("https-addr", ":443", "HTTPS network address")
		httpAddr    = fs.String("http-addr", ":80", "HTTP network address")
		cvPath      = fs.String("cv", "resume.yaml", "path to resume in YAML format")
		versionFlag = fs.Bool("version", false, "print version information and exit")
	)
	fs.Parse(os.Args[1:])

	if *versionFlag {
		fmt.Fprintf(os.Stdout, "Current build version %s\n", version)
		os.Exit(0)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	file, err := os.Open(*cvPath)
	if err != nil {
		errorLog.Fatal(err)
	}
	cv, err := yml.LoadCV(file)
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
		cv:            cv,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
	}

	go func() {
		_, tlsPort, err := net.SplitHostPort(*httpsAddr)
		if err != nil {
			errorLog.Fatal(err)
		}
		infoLog.Printf("Starting http server on %s", *httpAddr)
		httpSrv := http.Server{
			Addr:         *httpAddr,
			Handler:      app.recoverPanic(app.logRequest(app.httpsRedirect(tlsPort))),
			ErrorLog:     errorLog,
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		if err := httpSrv.ListenAndServe(); err != nil {
			errorLog.Fatal(err)
		}
	}()

	infoLog.Printf("Starting https server on %s", *httpsAddr)
	srv := http.Server{
		Addr:         *httpsAddr,
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServeTLS(*sslCert, *sslKey); err != nil {
		errorLog.Fatal(err)
	}
}
