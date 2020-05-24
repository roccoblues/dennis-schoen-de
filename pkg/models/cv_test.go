package models

import (
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestCV(t *testing.T) {
	data, err := ioutil.ReadFile("../../resume.yaml")
	if err != nil {
		t.Fatal(err)
	}
	cv := &CV{}
	if err = yaml.Unmarshal(data, &cv); err != nil {
		t.Fatal(err)
	}
}
