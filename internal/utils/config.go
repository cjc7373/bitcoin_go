package utils

import (
	"errors"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBPath string
}

func ParseConfig(configPath string) *Config {
	var conf Config
	data, err := os.ReadFile(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		conf.DBPath = "blockchain.db"

		newData, err := yaml.Marshal(&conf)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(configPath, newData, 0644)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	} else {
		err = yaml.Unmarshal(data, &conf)
		if err != nil {
			panic(err)
		}
	}

	return &conf
}
