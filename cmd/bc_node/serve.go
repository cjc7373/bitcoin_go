package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_server"
	"github.com/cjc7373/bitcoin_go/internal/utils"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
	"github.com/spf13/cobra"
)

var genesis bool
var connectTo string

func RunServe(cmd *cobra.Command, args []string) {
	config := utils.GetConfigFromContext(cmd.Context())
	w := wallet.ReadOrCreateWalletFromConfig(config)
	logger := slog.Default()

	service := network.NewService()
	rpcServer := rpc_server.NewRPCServer(service, logger)
	done := make(chan error)
	go rpc_server.Serve(rpcServer, config.ListenAddr, done)

	if genesis {
		fmt.Println(w)
	} else {
		rpcServer.RPCClient.ConnectFirstNode(connectTo, config.ListenAddr, config.NodeName)
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
	serveCmd.Flags().StringVar(&connectTo, "connect-to", "", "connect to a node as the first neighbour (either this flag or --genesis should be set)")
	serveCmd.MarkFlagsOneRequired("genesis", "connect-to")
	serveCmd.MarkFlagsMutuallyExclusive("genesis", "connect-to")
	return serveCmd
}
