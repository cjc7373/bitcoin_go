package network

import (
	"context"
	"log"

	"github.com/cjc7373/bitcoin_go/internal/network/proto"
)

func (s *Service) broadcastNodes(ttl uint32) {
	log.Println("broadcasting nodes..")
	s.RLock()
	log.Println("connected nodes: ", s.connectedNodes)
	protoNodes := make([]*proto.Node, 0, len(s.connectedNodes))
	clients := make([]proto.DiscoveryClient, 0, len(s.connectedNodes))
	for _, v := range s.connectedNodes {
		protoNodes = append(protoNodes, &proto.Node{Address: v.Address, Name: v.Name})
		clients = append(clients, v.DiscoveryClient)
	}
	s.RUnlock()
	for _, client := range clients {
		_, err := client.BroadcastNodes(context.Background(), &proto.NodeBroadcast{
			Nodes: protoNodes,
			TTL:   ttl,
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Service) ConnectFirstNode(remoteAddress, localAddress, localName string) error {
	_, err := s.connectNode(remoteAddress, "")
	if err != nil {
		return err
	}

	s.RLock()
	client := s.connectedNodes[remoteAddress].DiscoveryClient
	s.RUnlock()
	_, err = client.RequestNodes(context.Background(), &proto.Node{
		Name:    localName,
		Address: localAddress,
	})
	if err != nil {
		return err
	}
	return nil
}
