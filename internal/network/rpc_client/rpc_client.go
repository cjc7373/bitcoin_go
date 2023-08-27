package rpc_client

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
)

type statsHandler struct {
	service *network.Service
}

func (h *statsHandler) TagRPC(ctx context.Context, tagInfo *stats.RPCTagInfo) context.Context {
	return ctx
}

func (h *statsHandler) HandleRPC(context.Context, stats.RPCStats) {}

func (h *statsHandler) TagConn(ctx context.Context, tagInfo *stats.ConnTagInfo) context.Context {
	return ctx
}

func (h *statsHandler) HandleConn(ctx context.Context, connStats stats.ConnStats) {
	switch connStats.(type) {
	case *stats.ConnBegin:
		s := "unknown address"
		p, ok := peer.FromContext(ctx)
		if ok {
			s = p.Addr.String()
		}
		log.Println("connection established to ", s)
	case *stats.ConnEnd:
		p, ok := peer.FromContext(ctx)
		if !ok {
			log.Println(errors.New("unknown connection disconnected"))
		} else {
			addr := p.Addr.String()
			log.Printf("Connection with %v disconnected", addr)
			DisconnectNode(h.service, addr)
		}
	}
}

func ConnectNode(service *network.Service, address string, name string) (bool, error) {
	if _, ok := service.GetConnectedNode(address); ok {
		// already connected
		return false, nil
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithNoProxy(),
		grpc.WithStatsHandler(&statsHandler{service: service}),
	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return false, err
	}
	node := &network.Node{
		Node:          proto.Node{Name: name, Address: address},
		Conn:          conn,
		LastHeartbeat: time.Now(),
		BitcoinClient: proto.NewBitcoinClient(conn),
	}
	service.SetConnectedNode(address, node)
	return true, nil
}

func DisconnectNode(service *network.Service, address string) error {
	node, ok := service.GetConnectedNode(address)
	if !ok {
		return errors.New("node does not exist")
	}
	node.Conn.Close()
	service.DeleteConnectedNode(address)
	return nil
}

func ConnectFirstNode(service *network.Service, remoteServerAddr, localServerAddr, localName string) error {
	_, err := ConnectNode(service, remoteServerAddr, "")
	if err != nil {
		return err
	}

	node, ok := service.GetConnectedNode(remoteServerAddr)
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

func BroadcastNodes(service *network.Service, ttl uint32) {
	log.Println("broadcasting nodes..")
	connectedNodes := service.GetConnectedNodes()
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

func HandleBroadcast(service *network.Service) {
	timer := time.NewTimer(network.BroadcastInterval)
	for {
		select {
		case <-service.ShouldBroadcast:
			log.Println("shouldBroadcast chan received")
			if !timer.Stop() {
				<-timer.C
			}
		case <-timer.C:
			log.Println("broadcast timer expired")
		}
		BroadcastNodes(service, 1)
		timer.Reset(network.BroadcastInterval)
	}
}
