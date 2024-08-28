package rpc_client

import (
	"context"
	"log/slog"
	"time"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cjc7373/bitcoin_go/internal/block"
	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/proto"
)

type RPCClient struct {
	service *network.Service
	logger  *slog.Logger
	db      *bolt.DB
}

func NewRPCClient(service *network.Service, logger *slog.Logger, db *bolt.DB) RPCClient {
	return RPCClient{
		service: service,
		logger:  logger,
		db:      db,
	}
}

func (c *RPCClient) ConnectNode(address string, name string, localServerAddr, localName string) (bool, error) {
	if name == localName {
		return false, nil
	}
	if _, ok := c.service.GetConnectedNode(address); ok {
		// already connected
		return false, nil
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithNoProxy(),
		grpc.WithStatsHandler(&statsHandler{service: c.service, logger: c.logger}),
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
	c.logger.Info("node connected", "peer", node)
	node.BitcoinClient.SendNodes(context.TODO(), &proto.Nodes{Nodes: []*proto.Node{
		{
			Name:    localName,
			Address: localServerAddr,
		},
	}})
	c.logger.Info("send nodes rpc sent", "peer", node)
	return true, nil
}

func (c *RPCClient) getGrpcClients() []proto.BitcoinClient {
	connectedNodes := c.service.GetConnectedNodes()
	clients := make([]proto.BitcoinClient, 0, len(connectedNodes))
	for _, v := range connectedNodes {
		clients = append(clients, v.BitcoinClient)
	}
	return clients
}

// this function may block
func (c *RPCClient) BroadcastNodes(ttl uint32) {
	connectedNodes := c.service.GetConnectedNodes()
	c.logger.Info("broadcasting nodes..", "connectedNodes", connectedNodes)
	protoNodes := make([]*proto.Node, 0, len(connectedNodes))
	clients := make([]proto.BitcoinClient, 0, len(connectedNodes))
	for _, v := range connectedNodes {
		protoNodes = append(protoNodes, &proto.Node{Address: v.Address, Name: v.Name})
		clients = append(clients, v.BitcoinClient)
	}

	// TODO: use goroutine to parallel this
	for _, client := range clients {
		_, err := client.SendNodes(context.Background(), &proto.Nodes{
			Nodes: protoNodes,
		})
		if err != nil {
			c.logger.Error("", "err", err)
		}
	}
}

// this function may block
func (c *RPCClient) BroadcastBlockchain() {
	bc, err := block.GetBlockchain(c.db)
	if err != nil {
		c.logger.Error("", "err", err)
	}
	c.logger.Info("broadcasting blockchain", "blockchain", (*block_proto.Blockchain)(bc))
	for _, client := range c.getGrpcClients() {
		_, err := client.SendChainMetadata(context.Background(), bc)
		if err != nil {
			c.logger.Error("", "err", err)
		}
	}
}

func (c *RPCClient) HandleBroadcast() {
	timer := time.NewTimer(network.BroadcastInterval)
	for {
		select {
		case <-c.service.ShouldBroadcast:
			c.logger.Info("shouldBroadcast chan received")
			if !timer.Stop() {
				<-timer.C
			}
		case <-timer.C:
			c.logger.Info("broadcast timer expired")
		}
		c.BroadcastNodes(1)
		timer.Reset(network.BroadcastInterval)
	}
}
