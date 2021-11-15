package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var Conf *Config
var allowedModes = map[string]bool{
	"defence":        true,
	"attack-defence": true,
}

type Config struct {
	Mode            string `json:"mode"`
	Network         string `yaml:"network"`
	CheckerPassword string `yaml:"checker_password"`
	RoundInterval   string `yaml:"round_interval"`
	ExploitInterval string `yaml:"exploit_interval"`
	Teams           Teams  `yaml:"teams"`
	Users           Users  `yaml:"users"`
}

type Teams struct {
	Number    int       `yaml:"number"`
	Resources Resources `yaml:"resources"`
}

type Users struct {
	Number    int       `yaml:"number"`
	Resources Resources `yaml:"resources"`
}

type Resources struct {
	Memory int `yaml:"memory"`
	VCPU   int `yaml:"vcpu"`
}

func ReadConf(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	Conf = &Config{}
	err = yaml.Unmarshal(buf, Conf)
	if err != nil {
		return fmt.Errorf("in file %q: %v", filename, err)
	}
	if !allowedModes[Conf.Mode] {
		return fmt.Errorf("unsuported mode, allowed %v", allowedModes)
	}

	return nil
}
