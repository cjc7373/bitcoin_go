package rpc_server

import (
	"context"
	"log"
	"log/slog"
	"net"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/proto"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_client"
	"github.com/cjc7373/bitcoin_go/internal/utils"
	"google.golang.org/grpc"
)

type BitcoinServer struct {
	RPCClient rpc_client.RPCClient
	s         *network.Service
	logger    *slog.Logger
	config    *utils.Config

	proto.UnimplementedBitcoinServer
}

func NewRPCServer(service *network.Service, logger *slog.Logger, config *utils.Config) BitcoinServer {
	return BitcoinServer{
		s:         service,
		logger:    logger,
		RPCClient: rpc_client.NewRPCClient(service, logger),
		config:    config,
	}
}

func (d *BitcoinServer) SendNodes(ctx context.Context, nodes *proto.Nodes) (*proto.Empty, error) {
	connectedNodesUpdated := false
	for _, node := range nodes.Nodes {
		connected, err := d.RPCClient.ConnectNode(node.Address, node.Name, d.config.ListenAddr, d.config.NodeName)
		if err != nil {
			return nil, err
		}
		if connected {
			connectedNodesUpdated = true
		}
	}

	if connectedNodesUpdated {
		// don't block on this channel
		select {
		case d.s.ShouldBroadcast <- struct{}{}:
		default:
			d.logger.Info("ShouldBroadcast not ready to receive, skipping")
		}
	}

	return &proto.Empty{}, nil
}

func Serve(rpcServer BitcoinServer, address string, done chan error) {
	lis, _ := net.Listen("tcp", address)
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterBitcoinServer(grpcServer, &rpcServer)
	go rpcServer.RPCClient.HandleBroadcast()
	log.Println("serving at ", address)
	done <- grpcServer.Serve(lis)
}
