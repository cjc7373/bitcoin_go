package network

import (
	"context"

	"github.com/cjc7373/bitcoin_go/internal/network/proto"
)

func (s *service) broadcastNodes(ttl uint32) {
	s.RLock()
	protoNodes := make([]*proto.Node, len(s.connectedNodes))
	clients := make([]proto.DiscoveryClient, len(s.connectedNodes))
	for _, v := range s.connectedNodes {
		protoNodes = append(protoNodes, &proto.Node{Address: v.Address, Name: v.Name})
		clients = append(clients, v.DiscoveryClient)
	}
	s.RUnlock()

	for _, client := range clients {
		client.BroadcastNodes(context.Background(), &proto.NodeBroadcast{
			Nodes: protoNodes,
			TTL:   ttl,
		})
	}
}
