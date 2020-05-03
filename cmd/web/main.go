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
	hostName      string
	templateCache map[string]*template.Template
}

func main() {
	fs := flag.NewFlagSet("web", flag.ExitOnError)
	var (
		sslFlag     = fs.Bool("ssl", false, "enable HTTPS, requires --ssl-cert and --ssl-key, HTTP requests will be redirected")
		sslCert     = fs.String("ssl-cert", "", "SSL Certificate")
		sslKey      = fs.String("ssl-key", "", "SSL Key")
		httpsAddr   = fs.String("https-addr", ":443", "HTTPS network address")
		httpAddr    = fs.String("http-addr", ":80", "HTTP network address")
		hostName    = fs.String("hostname", "", "redirect unknown host requests to hostname (optional)")
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

	if *sslFlag {
		if *sslCert == "" {
			errorLog.Fatal("--ssl-cert is required")
		}
		if *sslKey == "" {
			errorLog.Fatal("--ssl-key is required")
		}
	}

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
		hostName:      *hostName,
	}

	_, tlsPort, err := net.SplitHostPort(*httpsAddr)
	if err != nil {
		errorLog.Fatal(err)
	}

	var g run.Group
	{
		if *sslFlag {
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
			var httpHandler http.Handler
			if *sslFlag {
				// redirect HTTP requests to HTTPS if SSL is enabled
				httpHandler = app.recoverPanic(app.logRequest(app.httpsRedirect(tlsPort)))
			} else {
				httpHandler = app.routes()
			}
			server := http.Server{
				Addr:         *httpAddr,
				Handler:      httpHandler,
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
			g.Add(run.SignalHandler(context.Background(), os.Interrupt))
		}
	}

	infoLog.Printf("Exit with: %s", g.Run())
}
