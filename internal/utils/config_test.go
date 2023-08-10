package utils

import (
	"fmt"
	"testing"

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
