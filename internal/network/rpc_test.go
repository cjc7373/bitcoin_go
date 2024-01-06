package network_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_client"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_server"
	"github.com/stretchr/testify/assert"
)

func TestDiscovery(t *testing.T) {
	serverAddr1 := "127.0.0.1:12200"
	service1 := network.NewService()
	done := make(chan error)
	go func() {
		rpc_server.Serve(service1, serverAddr1, done)
	}()

	serverAddr2 := "127.0.0.1:12201"
	service2 := network.NewService()
	go func() {
		rpc_server.Serve(service2, serverAddr2, done)
	}()

	// wait server start
	time.Sleep(time.Microsecond * 100)

	err := rpc_client.ConnectFirstNode(service2, serverAddr1, serverAddr2, "service2")
	assert.Nil(t, err)

	s2Nodes := service2.GetConnectedNodes()
	fmt.Println(s2Nodes)
	assert.Nil(t, err)
	assert.Len(t, s2Nodes, 1)

	rpc_client.DisconnectNode(service2, serverAddr1)

	// wait server handle conn close
	time.Sleep(time.Microsecond * 1000)
	assert.Len(t, service2.GetConnectedNodes(), 0)
}
