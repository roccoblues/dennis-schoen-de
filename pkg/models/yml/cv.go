package yml

import (
	"io"

	"github.com/roccoblues/dennis-schoen.de/pkg/models"
	"gopkg.in/yaml.v2"
)

func LoadCV(r io.Reader) (*models.CV, error) {
	cv := &models.CV{}

	d := yaml.NewDecoder(r)
	err := d.Decode(&cv)
	if err != nil {
		return nil, err
	}

	return cv, nil
}
