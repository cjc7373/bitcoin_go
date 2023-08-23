package network

import (
	"sync"
	"time"

	"github.com/cjc7373/bitcoin_go/internal/network/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Node struct {
	proto.Node

	Conn            *grpc.ClientConn
	LastHeartbeat   time.Time
	DiscoveryClient proto.DiscoveryClient
}

type service struct {
	sync.RWMutex
	connectedNodes map[string]*Node
}

func NewService() *service {
	return &service{
		connectedNodes: make(map[string]*Node),
	}
}

// lock must not hold when calling this method
func (s *service) connectNode(address string, name string) (bool, error) {
	s.RLock()
	_, ok := s.connectedNodes[name]
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
	s.connectedNodes[name] = &Node{
		Node:            proto.Node{Name: name, Address: address},
		Conn:            conn,
		LastHeartbeat:   time.Now(),
		DiscoveryClient: proto.NewDiscoveryClient(conn),
	}
	return true, nil
}
