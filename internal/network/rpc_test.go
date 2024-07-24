package network_test

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_client"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_server"
	"github.com/cjc7373/bitcoin_go/internal/utils"
	"github.com/stretchr/testify/assert"
)

// TODO: refactor with ginkgo
func TestDiscovery(t *testing.T) {
	logger := slog.Default()

	config1 := utils.Config{
		ListenAddr: "127.0.0.1:12201",
		NodeName:   "node1",
	}
	service1 := network.NewService()
	done := make(chan error)
	rpcServer1 := rpc_server.NewRPCServer(service1, logger.With("node", config1.NodeName), &config1)
	go func() {
		rpc_server.Serve(rpcServer1, config1.ListenAddr, done)
	}()

	config2 := utils.Config{
		ListenAddr: "127.0.0.1:12202",
		NodeName:   "node2",
	}
	service2 := network.NewService()
	rpcServer2 := rpc_server.NewRPCServer(service2, logger.With("node", config2.NodeName), &config2)
	go func() {
		rpc_server.Serve(rpcServer2, config2.ListenAddr, done)
	}()

	config3 := utils.Config{
		ListenAddr: "127.0.0.1:12203",
		NodeName:   "node3",
	}
	service3 := network.NewService()
	rpcServer3 := rpc_server.NewRPCServer(service3, logger.With("node", config3.NodeName), &config3)
	go func() {
		rpc_server.Serve(rpcServer3, config3.ListenAddr, done)
	}()

	// wait server start
	time.Sleep(time.Microsecond * 100)

	connected, err := rpcServer2.RPCClient.ConnectNode(config1.ListenAddr, config1.NodeName, config2.ListenAddr, config2.NodeName)
	assert.Nil(t, err)
	assert.True(t, connected)

	s2Nodes := service2.GetConnectedNodes()
	fmt.Println(s2Nodes)
	assert.Nil(t, err)
	assert.Len(t, s2Nodes, 1)

	rpc_client.DisconnectNode(service2, config1.ListenAddr)

	// wait server handle conn close
	time.Sleep(time.Microsecond * 1000)
	assert.Len(t, service2.GetConnectedNodes(), 0)
}
