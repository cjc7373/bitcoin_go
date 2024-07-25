package network_test

import (
	"fmt"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_server"
	"github.com/cjc7373/bitcoin_go/internal/utils"
)

var _ = Describe("RPC test", func() {
	var rpcServers []*rpc_server.BitcoinServer
	var configs []utils.Config
	var nodeNum = 5
	BeforeEach(func() {
		rpcServers = make([]*rpc_server.BitcoinServer, 0)
		configs = make([]utils.Config, 0)
		logger := slog.Default()
		for i := range nodeNum {
			config := utils.Config{
				ListenAddr: fmt.Sprintf("127.0.0.1:1220%v", i),
				NodeName:   fmt.Sprintf("node%v", i),
			}
			configs = append(configs, config)
			service := network.NewService()
			rpcServer := rpc_server.NewRPCServer(service, logger.With("node", config.NodeName), &config)
			rpcServer.Serve()
			rpcServers = append(rpcServers, &rpcServer)
		}
	})

	AfterEach(func() {
		for _, rpcServer := range rpcServers {
			rpcServer.Stop()
		}
		rpcServers = nil
		configs = nil
	})

	It("connects a node", func() {
		connected, err := rpcServers[0].RPCClient.ConnectNode(
			configs[1].ListenAddr, configs[1].NodeName, configs[0].ListenAddr, configs[0].NodeName,
		)
		Expect(connected).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() int {
			return len(rpcServers[0].GetConnectedNodes())
		}).Should(Equal(1))

		Expect(rpcServers[0].DisconnectNode(configs[1].ListenAddr)).Should(Succeed())
		Eventually(func() int {
			return len(rpcServers[0].GetConnectedNodes())
		}).Should(Equal(0))
	})

	It("connects all nodes", func() {
		for i := 1; i < nodeNum; i++ {
			connected, err := rpcServers[0].RPCClient.ConnectNode(
				configs[i].ListenAddr, configs[i].NodeName, configs[0].ListenAddr, configs[0].NodeName,
			)
			Expect(connected).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		}

		for i := range nodeNum {
			Eventually(func() int {
				return len(rpcServers[i].GetConnectedNodes())
			}).Should(Equal(4))
		}
	})
})
