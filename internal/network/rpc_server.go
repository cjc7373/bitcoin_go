package network

import (
	"context"
	"errors"

	"github.com/cjc7373/bitcoin_go/internal/network/proto"
	"google.golang.org/grpc/peer"
)

type discoveryServer struct {
	s *service

	proto.UnimplementedDiscoveryServer
}

func (d *discoveryServer) RequestNodes(ctx context.Context, nodeRequest *proto.NodeRequest) (*proto.Empty, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("cannot get peer info")
	}
	d.s.connectNode(p.Addr.String(), nodeRequest.Name)
	return &proto.Empty{}, nil
}

func (d *discoveryServer) BroadcastNodes(ctx context.Context, nodeBroadcast *proto.NodeBroadcast) (*proto.Empty, error) {
	for _, node := range nodeBroadcast.Nodes {
		// TODO: exclude non-modified nodes and sender node
		_, err := d.s.connectNode(node.Address, node.Name)
		if err != nil {
			return nil, err
		}
	}

	if nodeBroadcast.TTL > 0 {
		d.s.broadcastNodes(nodeBroadcast.TTL - 1)
	}

	return &proto.Empty{}, nil
}
