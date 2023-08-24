package utils

import (
	"context"
	"errors"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBPath  string
	Wallets map[string]string
}

func NewDefaultConfig() *Config {
	return &Config{
		DBPath:  "blockchain.db",
		Wallets: map[string]string{},
	}
}

func ParseConfig(configPath string) *Config {
	var conf Config
	data, err := os.ReadFile(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		conf = *NewDefaultConfig()
		conf.WriteToFile(configPath)
	} else if err != nil {
		panic(err)
	} else {
		err = yaml.Unmarshal(data, &conf)
		if err != nil {
			panic(err)
		}
	}

	if conf.Wallets == nil {
		conf.Wallets = make(map[string]string)
	}

	return &conf
}

func (conf *Config) WriteToFile(configPath string) {
	newData, err := yaml.Marshal(&conf)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(configPath, newData, 0644)
	if err != nil {
		panic(err)
	}
}

var ConfigKey struct{}

func GetConfigFromContext(ctx context.Context) *Config {
	config, ok := ctx.Value(&ConfigKey).(*Config)
	if !ok {
		panic("cannot get config from context")
	}
	return config
}
