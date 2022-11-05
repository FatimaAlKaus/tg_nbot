package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BotToken string `yaml:"token"`
}

func Load(file string) (*Config, error) {

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var c Config
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	return &c, err
}
