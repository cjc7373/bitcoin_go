package utils

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gopkg.in/yaml.v3"
)

var _ = Describe("config test", func() {
	It("tests yaml", func() {
		type T struct {
			A string
			B []byte
			C map[string]string

			unexported string
		}

		data := T{"xxx", []byte("xaaaa"), map[string]string{"key": "v1", "k3": "v3"}, "foo"}
		out, _ := yaml.Marshal(data)
		fmt.Println(out)
		var m T
		err := yaml.Unmarshal(out, &m)
		fmt.Println(err, m)
	})

	It("parses config", func() {
		dataDir := "./testdata"
		config := ParseConfig(dataDir)

		Expect(config.DBPath).To(Equal("blockchain.db"))
		Expect(strings.TrimSpace(config.Wallets["default"])).To(Equal("foo"))
		Expect(config.dataDir).To(Equal(dataDir))
	})
})
