package network

import (
	"context"
	"log"

	"github.com/cjc7373/bitcoin_go/internal/network/proto"
)

type discoveryServer struct {
	s *Service

	proto.UnimplementedDiscoveryServer
}

func (d *discoveryServer) RequestNodes(ctx context.Context, nodeRequest *proto.Node) (*proto.Empty, error) {
	_, err := d.s.connectNode(nodeRequest.Address, nodeRequest.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	d.s.shouldBroadcast <- struct{}{}
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
