package rpc_server

import (
	"context"
	"log"
	"log/slog"
	"net"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/proto"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_client"
	"google.golang.org/grpc"
)

type BitcoinServer struct {
	RPCClient rpc_client.RPCClient
	s         *network.Service
	logger    *slog.Logger

	proto.UnimplementedBitcoinServer
}

func NewRPCServer(service *network.Service, logger *slog.Logger) BitcoinServer {
	return BitcoinServer{
		s:         service,
		logger:    logger,
		RPCClient: rpc_client.NewRPCClient(service, logger),
	}
}

func (d *BitcoinServer) RequestNodes(ctx context.Context, nodeRequest *proto.Node) (*proto.Empty, error) {
	_, err := d.RPCClient.ConnectNode(nodeRequest.Address, nodeRequest.Name)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}
	// don't block on this channel
	select {
	case d.s.ShouldBroadcast <- struct{}{}:
	default:
		d.logger.Info("ShouldBroadcast not ready to receive, skipping")
	}
	return &proto.Empty{}, nil
}

func (d *BitcoinServer) BroadcastNodes(ctx context.Context, nodeBroadcast *proto.NodeBroadcast) (*proto.Empty, error) {
	for _, node := range nodeBroadcast.Nodes {
		// TODO: exclude non-modified nodes and sender node
		_, err := d.RPCClient.ConnectNode(node.Address, node.Name)
		if err != nil {
			return nil, err
		}
	}

	if nodeBroadcast.TTL > 0 {
		d.RPCClient.BroadcastNodes(nodeBroadcast.TTL - 1)
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
