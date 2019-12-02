package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hashicorp/hcl"
	"github.com/oklog/run"
	"github.com/roccoblues/dennis-schoen.de/pkg/models"
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
		sslCert     = fs.String("ssl-cert", "", "SSL Certificate")
		sslKey      = fs.String("ssl-key", "", "SSL Key")
		httpsAddr   = fs.String("https-addr", ":443", "HTTPS network address")
		httpAddr    = fs.String("http-addr", ":80", "HTTP network address")
		cvPath      = fs.String("cv", "resume.conf", "path to resume in HCL format")
		versionFlag = fs.Bool("version", false, "print version information and exit")
	)
	fs.Parse(os.Args[1:])

	if *versionFlag {
		fmt.Fprintf(os.Stdout, "Current build version %s\n", version)
		os.Exit(0)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var cv *models.CV
	{
		hclCV, err := ioutil.ReadFile(*cvPath)
		if err != nil {
			errorLog.Fatal(err)
		}
		cv = &models.CV{}
		if err = hcl.Unmarshal(hclCV, &cv); err != nil {
			errorLog.Fatal(err)
		}
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

	_, tlsPort, err := net.SplitHostPort(*httpsAddr)
	if err != nil {
		errorLog.Fatal(err)
	}

	var g run.Group
	{
		{
			infoLog.Printf("Starting HTTPS server on %s", *httpsAddr)
			server := http.Server{
				Addr:     *httpsAddr,
				Handler:  app.routes(),
				ErrorLog: errorLog,
				TLSConfig: &tls.Config{
					PreferServerCipherSuites: true,
				},
				IdleTimeout:  time.Minute,
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}
			g.Add(func() error {
				return server.ListenAndServeTLS(*sslCert, *sslKey)
			}, func(error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				infoLog.Printf("Shutting down HTTPS server")
				server.Shutdown(ctx)
			})
		}
		{
			infoLog.Printf("Starting HTTP server on %s", *httpAddr)
			server := http.Server{
				Addr:         *httpAddr,
				Handler:      app.recoverPanic(app.logRequest(app.httpsRedirect(tlsPort))),
				ErrorLog:     errorLog,
				IdleTimeout:  time.Minute,
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}
			g.Add(func() error {
				return server.ListenAndServe()
			}, func(error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				infoLog.Printf("Shutting down HTTP server")
				server.Shutdown(ctx)
			})
		}
		{
			// Catch ctrl-C
			var (
				ctx, cancel = context.WithCancel(context.Background())
				sigchan     = make(chan os.Signal, 1)
			)
			signal.Notify(sigchan, os.Interrupt)
			g.Add(func() error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case sig := <-sigchan:
					return fmt.Errorf("received signal %s", sig)
				}
			}, func(error) {
				cancel()
			})
		}
	}

	infoLog.Printf("Exit with: %s", g.Run())
}
