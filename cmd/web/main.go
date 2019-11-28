package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
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
		addr        = fs.String("addr", ":80", "HTTP network address")
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

	infoLog.Printf("Starting server on %s", *addr)
	srv := http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}
