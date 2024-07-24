package network_test

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_client"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_server"
	"github.com/stretchr/testify/assert"
)

func TestDiscovery(t *testing.T) {
	logger := slog.Default()

	serverAddr1 := "127.0.0.1:12200"
	service1 := network.NewService()
	done := make(chan error)
	rpcServer1 := rpc_server.NewRPCServer(service1, logger.With("server", 1))
	go func() {
		rpc_server.Serve(rpcServer1, serverAddr1, done)
	}()

	serverAddr2 := "127.0.0.1:12201"
	service2 := network.NewService()
	rpcServer2 := rpc_server.NewRPCServer(service2, logger.With("server", 2))
	go func() {
		rpc_server.Serve(rpcServer2, serverAddr2, done)
	}()

	serverAddr3 := "127.0.0.1:12202"
	service3 := network.NewService()
	rpcServer3 := rpc_server.NewRPCServer(service3, logger.With("server", 3))
	go func() {
		rpc_server.Serve(rpcServer3, serverAddr3, done)
	}()

	// wait server start
	time.Sleep(time.Microsecond * 100)

	err := rpcServer2.RPCClient.ConnectFirstNode(serverAddr1, serverAddr2, "service2")
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
