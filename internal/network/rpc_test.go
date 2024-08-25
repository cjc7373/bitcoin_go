package network_test

import (
	"fmt"
	"log/slog"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cjc7373/bitcoin_go/internal/network/rpc_server"
	"github.com/cjc7373/bitcoin_go/internal/utils"
)

var _ = Describe("RPC test", func() {
	var rpcServers []*rpc_server.BitcoinServer
	var configs []*utils.Config
	var nodeNum = 5
	BeforeEach(func() {
		rpcServers = make([]*rpc_server.BitcoinServer, 0)
		configs = make([]*utils.Config, 0)
		logger := slog.Default()
		for i := range nodeNum {
			config := utils.ParseConfig(fmt.Sprintf("./testdata/%v", i))
			config.ListenAddr = fmt.Sprintf("127.0.0.1:1220%v", i)
			config.NodeName = fmt.Sprintf("node%v", i)

			configs = append(configs, config)
			rpcServer := rpc_server.NewRPCServer(logger.With("node", config.NodeName), config)
			rpcServer.Serve()
			rpcServers = append(rpcServers, &rpcServer)
		}
	})

	AfterEach(func() {
		for _, rpcServer := range rpcServers {
			rpcServer.Stop()
		}
		for i := range nodeNum {
			if err := os.RemoveAll(fmt.Sprintf("testdata/%v", i)); err != nil {
				panic(err)
			}
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
			connected, err := rpcServers[i].RPCClient.ConnectNode(
				configs[0].ListenAddr, configs[0].NodeName, configs[i].ListenAddr, configs[i].NodeName,
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
