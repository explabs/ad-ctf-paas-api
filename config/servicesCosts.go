package config

import (
	"fmt"
	"github.com/explabs/ad-ctf-paas-api/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ServicesInfo struct {
	Services []*models.Service `yaml:"services"`
}

func (s *ServicesInfo) Load() error {
	buf, err := ioutil.ReadFile("checker.yml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, &s)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}
