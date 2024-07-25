package network

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestRPC(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RPC Suite")
}
