package models

import (
	"html/template"
)

type CV struct {
	Name    string
	Title   string
	Address struct {
		Street string
		Place  string
	}
	Phone      string
	Email      string
	Homepage   string
	Birthday   string
	Summary    []string
	Experience []struct {
		Title       string
		StartDate   string
		EndDate     string
		Description template.HTML
		Technology  []string
		Employer    struct {
			Name     string
			Location string
			URL      template.URL
		}
	}
	Languages []struct {
		Language string
		Level    string
	}
	Trainings []struct {
		Title  string
		Year   int
		Issuer struct {
			Name     string
			URL      template.URL
			Location string
		}
	}
}
