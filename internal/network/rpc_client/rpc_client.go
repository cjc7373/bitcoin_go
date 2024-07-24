package rpc_client

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"time"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RPCClient struct {
	service *network.Service
	logger  *slog.Logger
}

func NewRPCClient(service *network.Service, logger *slog.Logger) RPCClient {
	return RPCClient{
		service: service,
		logger:  logger,
	}
}

func (c *RPCClient) ConnectNode(address string, name string) (bool, error) {
	if _, ok := c.service.GetConnectedNode(address); ok {
		// already connected
		return false, nil
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithNoProxy(),
		grpc.WithStatsHandler(&statsHandler{service: c.service}),
	}
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return false, err
	}
	node := &network.Node{
		Node:          proto.Node{Name: name, Address: address},
		Conn:          conn,
		LastHeartbeat: time.Now(),
		BitcoinClient: proto.NewBitcoinClient(conn),
	}
	c.service.SetConnectedNode(address, node)
	return true, nil
}

func (c *RPCClient) ConnectFirstNode(remoteServerAddr, localServerAddr, localName string) error {
	_, err := c.ConnectNode(remoteServerAddr, "")
	if err != nil {
		return err
	}

	node, ok := c.service.GetConnectedNode(remoteServerAddr)
	if !ok {
		return errors.New("cannot get connected node")
	}
	_, err = node.BitcoinClient.RequestNodes(context.Background(), &proto.Node{
		Name:    localName,
		Address: localServerAddr,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *RPCClient) BroadcastNodes(ttl uint32) {
	log.Println("broadcasting nodes..")
	connectedNodes := c.service.GetConnectedNodes()
	log.Println("connected nodes: ", connectedNodes)
	protoNodes := make([]*proto.Node, 0, len(connectedNodes))
	clients := make([]proto.BitcoinClient, 0, len(connectedNodes))
	for _, v := range connectedNodes {
		protoNodes = append(protoNodes, &proto.Node{Address: v.Address, Name: v.Name})
		clients = append(clients, v.BitcoinClient)
	}
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

func (c *RPCClient) HandleBroadcast() {
	timer := time.NewTimer(network.BroadcastInterval)
	for {
		select {
		case <-c.service.ShouldBroadcast:
			log.Println("shouldBroadcast chan received")
			if !timer.Stop() {
				<-timer.C
			}
		case <-timer.C:
			log.Println("broadcast timer expired")
		}
		c.BroadcastNodes(1)
		timer.Reset(network.BroadcastInterval)
	}
}
