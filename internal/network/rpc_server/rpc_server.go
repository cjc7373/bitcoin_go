package rpc_server

import (
	"context"
	"errors"
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

	// whether server has exited
	done       chan error
	grpcServer *grpc.Server

	proto.UnimplementedBitcoinServer
}

func NewRPCServer(service *network.Service, logger *slog.Logger, config *utils.Config) BitcoinServer {
	return BitcoinServer{
		s:         service,
		logger:    logger,
		RPCClient: rpc_client.NewRPCClient(service, logger),
		config:    config,

		done: make(chan error),
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

func (d *BitcoinServer) Serve() error {
	lis, err := net.Listen("tcp", d.config.ListenAddr)
	if err != nil {
		return errors.Join(errors.New("listen failed"), err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	d.grpcServer = grpcServer
	proto.RegisterBitcoinServer(grpcServer, d)
	go d.RPCClient.HandleBroadcast()
	go func() {
		d.logger.Info("serving", "addr", d.config.ListenAddr)
		d.done <- grpcServer.Serve(lis)
	}()
	return nil
}

func (d *BitcoinServer) Stop() {
	d.grpcServer.GracefulStop()
}

func (d *BitcoinServer) GetConnectedNodes() map[string]*network.Node {
	return d.s.GetConnectedNodes()
}

func (d *BitcoinServer) DisconnectNode(address string) error {
	return rpc_client.DisconnectNode(d.s, address)
}
