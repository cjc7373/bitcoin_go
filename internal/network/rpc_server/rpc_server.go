package rpc_server

import (
	"context"
	"log"
	"net"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/proto"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_client"
	"google.golang.org/grpc"
)

type BitcoinServer struct {
	s *network.Service

	proto.UnimplementedBitcoinServer
}

func (d *BitcoinServer) RequestNodes(ctx context.Context, nodeRequest *proto.Node) (*proto.Empty, error) {
	_, err := rpc_client.ConnectNode(d.s, nodeRequest.Address, nodeRequest.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	d.s.ShouldBroadcast <- struct{}{}
	return &proto.Empty{}, nil
}

func (d *BitcoinServer) BroadcastNodes(ctx context.Context, nodeBroadcast *proto.NodeBroadcast) (*proto.Empty, error) {
	for _, node := range nodeBroadcast.Nodes {
		// TODO: exclude non-modified nodes and sender node
		_, err := rpc_client.ConnectNode(d.s, node.Address, node.Name)
		if err != nil {
			return nil, err
		}
	}

	if nodeBroadcast.TTL > 0 {
		rpc_client.BroadcastNodes(d.s, nodeBroadcast.TTL-1)
	}

	return &proto.Empty{}, nil
}

func Serve(service *network.Service, address string, done chan error) {
	lis, _ := net.Listen("tcp", address)
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	discovery := BitcoinServer{s: service}
	proto.RegisterBitcoinServer(grpcServer, &discovery)
	go rpc_client.HandleBroadcast(service)
	log.Println("serving at ", address)
	done <- grpcServer.Serve(lis)
}
