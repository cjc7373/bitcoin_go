package network

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/cjc7373/bitcoin_go/internal/network/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
)

type Node struct {
	proto.Node

	Conn            *grpc.ClientConn
	LastHeartbeat   time.Time
	DiscoveryClient proto.DiscoveryClient
}

func (node *Node) String() string {
	return fmt.Sprintf("Name: %v", node.Node.Name)
}

const broadcastInterval = time.Second * 30

type Service struct {
	shouldBroadcast chan struct{}

	// below fields are protected by RW lock
	sync.RWMutex
	// key is peer's address
	// contains connections initiated as a client
	connectedNodes map[string]*Node
}

func NewService() *Service {
	return &Service{
		shouldBroadcast: make(chan struct{}),
		connectedNodes:  make(map[string]*Node),
	}
}

// lock must not hold when calling this method
func (s *Service) connectNode(address string, name string) (bool, error) {
	s.RLock()
	_, ok := s.connectedNodes[address]
	s.RUnlock()
	if ok {
		// already connected
		return false, nil
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithNoProxy()}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return false, err
	}
	s.Lock()
	defer s.Unlock()
	s.connectedNodes[address] = &Node{
		Node:            proto.Node{Name: name, Address: address},
		Conn:            conn,
		LastHeartbeat:   time.Now(),
		DiscoveryClient: proto.NewDiscoveryClient(conn),
	}
	return true, nil
}

// lock must not hold when calling this method
func (s *Service) disconnectNode(address string) error {
	s.Lock()
	defer s.Unlock()
	node, ok := s.connectedNodes[address]
	if !ok {
		return fmt.Errorf("node %v does not exist in connectedNodes map", address)
	}
	node.Conn.Close()
	delete(s.connectedNodes, address)
	return nil
}

type statsHandler struct {
	service *Service
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
		log.Println("connection established from ", s)
	case *stats.ConnEnd:
		log.Println("Connection end", connStats, ctx)
		p, ok := peer.FromContext(ctx)
		if !ok {
			log.Println(errors.New("unknown client disconnected"))
		} else {
			err := h.service.disconnectNode(p.Addr.String())
			if err != nil {
				log.Println(err)
			} else {
				log.Println("connection ended from ", p.Addr.String())
			}
		}
	}
}

func (s *Service) handleBroadcast() {
	timer := time.NewTimer(broadcastInterval)
	for {
		select {
		case <-s.shouldBroadcast:
			log.Println("shouldBroadcast chan received")
			if !timer.Stop() {
				<-timer.C
			}
		case <-timer.C:
			log.Println("broadcast timer expired")
		}
		s.broadcastNodes(1)
		timer.Reset(broadcastInterval)
	}
}

func (s *Service) Serve(address string, done chan error) {
	lis, _ := net.Listen("tcp", address)
	opts := []grpc.ServerOption{grpc.StatsHandler(&statsHandler{service: s})}
	grpcServer := grpc.NewServer(opts...)
	discovery := discoveryServer{s: s}
	proto.RegisterDiscoveryServer(grpcServer, &discovery)
	go s.handleBroadcast()
	done <- grpcServer.Serve(lis)
}
