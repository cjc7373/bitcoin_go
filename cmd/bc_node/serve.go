package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_server"
	"github.com/cjc7373/bitcoin_go/internal/utils"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
)

var genesis bool
var connectToAddr string
var connectToNodeName string

func RunServe(cmd *cobra.Command, args []string) {
	config := utils.GetConfigFromContext(cmd.Context())
	w := wallet.ReadOrCreateWalletFromConfig(config)
	logger := slog.Default()

	service := network.NewService()
	rpcServer := rpc_server.NewRPCServer(service, logger, config)
	done := make(chan error)
	go rpcServer.Serve()

	if genesis {
		fmt.Println(w)
	} else {
		rpcServer.RPCClient.ConnectNode(connectToAddr, connectToNodeName, config.ListenAddr, config.NodeName)
	}
	log.Println(<-done)
}

func NewCmdServe() *cobra.Command {
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run a node",
		Run:   RunServe,
	}
	serveCmd.Flags().BoolVar(&genesis, "genesis", false, "if there's no blockchain exist,  create the genesis block")
	serveCmd.Flags().StringVar(&connectToAddr, "connect-to-addr", "", "address of a node to which connects as the first neighbour (either this flag or --genesis should be set)")
	serveCmd.Flags().StringVar(&connectToNodeName, "connect-to-node-name", "", "name of a node to which connects (should be used combined with connect-to-addr)")
	serveCmd.MarkFlagsOneRequired("genesis", "connect-to-addr")
	serveCmd.MarkFlagsMutuallyExclusive("genesis", "connect-to")
	serveCmd.MarkFlagsRequiredTogether("connect-to-addr", "connect-to-node-name")
	return serveCmd
}
