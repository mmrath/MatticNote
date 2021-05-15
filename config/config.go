package config

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

var Config MNConfig

func LoadConf() error {
	file, err := ioutil.ReadFile("matticnote.toml")
	if err != nil {
		return err
	}

	err = toml.Unmarshal(file, &Config)
	if err != nil {
		return err
	}

	return nil
}
