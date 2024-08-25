package utils

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBPath     string
	ListenAddr string
	NodeName   string
	Wallets    map[string]string

	dataDir string // keep this field private to avoid writing it to config file
}

func (c *Config) GetDataDir() string {
	return c.dataDir
}

func NewDefaultConfig() *Config {
	return &Config{
		DBPath:     "blockchain.db",
		ListenAddr: ":12000",
		NodeName:   "defaultNode",
		Wallets:    map[string]string{},
	}
}

func ParseConfig(dataDir string) *Config {
	configPath := path.Join(dataDir, "config.yaml")
	var conf Config
	data, err := os.ReadFile(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			panic(err)
		}
		conf = *NewDefaultConfig()
		conf.writeToFile(configPath)
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
	conf.dataDir = dataDir

	return &conf
}

// write config back to file
func (conf *Config) Write() {
	conf.writeToFile(conf.dataDir)
}

func (conf *Config) writeToFile(configPath string) {
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
