package models

import "html/template"

type CV struct {
	Summary    []string
	Experience []*WorkExperience
	Languages  []*struct {
		Language string
		Level    string
	}
	Trainings []*Training
}

type WorkExperience struct {
	Title       string
	StartDate   string `yaml:"startDate"`
	EndDate     string `yaml:"endDate"`
	Employer    *Employer
	Description template.HTML
}

type Training struct {
	Title  string
	Year   int
	Issuer struct {
		Name     string
		URL      template.URL
		Location string
	}
}

type Employer struct {
	Name     string
	Location string
	URL      template.URL
}
