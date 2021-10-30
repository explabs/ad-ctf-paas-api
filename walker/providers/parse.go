package providers

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfigProviders struct {
	Teams   []Team    `yaml:"teams"`
	Service []Service `yaml:"services"`
}

type Service struct {
	Name  string   `yaml:"name"`
	Put   []Script `yaml:"put"`
	Check []Script `yaml:"check"`
}

type Team struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	IP     string `yaml:"ip"`
}

type Script struct {
	Name string `yaml:"name"`
}

func (p *ConfigProviders) Parse(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, &p)
	if err != nil {
		return fmt.Errorf("in file %q: %v", filename, err)
	}
	return nil
}

type RoundsStruct struct {
	Rounds []Round `yaml:"rounds"`
}
type Round struct {
	Exploits []Exploit `yaml:"exploits"`
	News     string    `yaml:"news"`
	HintNews string    `yaml:"hint_news"`
}
type Exploit struct {
	ServiceName string `yaml:"service_name"`
	ScriptName  string `yaml:"script_name"`
}

func (r *RoundsStruct) Parse(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, &r)
	if err != nil {
		return fmt.Errorf("in file %q: %v", filename, err)
	}
	return nil
}
