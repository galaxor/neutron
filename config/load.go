package config

import (
	"encoding/json"
	"io/ioutil"
)

func Load(path string) (config *Config, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &config)
	return
}
