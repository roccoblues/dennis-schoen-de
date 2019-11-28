package main

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/markbates/pkger"
	"github.com/roccoblues/dennis-schoen.de/pkg/models"
)

type templateData struct {
	CV *models.CV
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// read all layouts into a buffer
	layouts := new(bytes.Buffer)
	err := pkger.Walk(dir, func(path string, info os.FileInfo, err error) error {
		matched, err := filepath.Match("*.layout.tmpl", info.Name())
		if err != nil {
			return err
		}
		if !matched {
			return nil
		}

		file, err := pkger.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = layouts.ReadFrom(file)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// now compile all pages together with the layouts
	err = pkger.Walk(dir, func(path string, info os.FileInfo, err error) error {
		matched, err := filepath.Match("*.page.tmpl", info.Name())
		if err != nil {
			return err
		}
		if !matched {
			return nil
		}

		file, err := pkger.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(file)

		t, err := template.New(info.Name()).Parse(buf.String() + layouts.String())
		if err != nil {
			return err
		}

		cache[info.Name()] = t
		return nil
	})

	if err != nil {
		return nil, err
	}

	return cache, nil
}
