package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestYaml(t *testing.T) {
	type T struct {
		A string
		B []byte
		C map[string]string
	}

	data := T{"xxx", []byte("xaaaa"), map[string]string{"key": "v1", "k3": "v3"}}
	out, _ := yaml.Marshal(data)
	fmt.Println(out)
	var m T
	err := yaml.Unmarshal([]byte{}, &m)
	fmt.Println(err, m)
}

func TestParseConfig(t *testing.T) {
	assert := assert.New(t)
	config := ParseConfig("./testdata/test_conf.yaml")

	assert.Equal("blockchain.db", config.DBPath)
	assert.Equal("foo", strings.TrimSpace(config.Wallets["default"]))
}
