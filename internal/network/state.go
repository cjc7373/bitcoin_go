package network

import (
	"fmt"
	"maps"
	"sync"
	"time"

	"google.golang.org/grpc"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/common"
	"github.com/cjc7373/bitcoin_go/internal/network/proto"
)

type Node struct {
	proto.Node

	Conn          *grpc.ClientConn
	LastHeartbeat time.Time
	BitcoinClient proto.BitcoinClient
}

func (node *Node) String() string {
	return fmt.Sprintf("Name: %v", node.Node.Name)
}

const BroadcastInterval = time.Second * 30

type Service struct {
	ShouldBroadcast chan struct{}

	// below fields are protected by RW lock
	lock sync.RWMutex
	// key is peer's server address
	// contains connections initiated as a client
	connectedNodes map[string]*Node

	// since we need to check if a tx is already in the queue, use a map
	pendingTxQ map[common.TxHashSized]*block_proto.Transaction
}

func (s *Service) GetConnectedNode(serverAddr string) (*Node, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	node, ok := s.connectedNodes[serverAddr]
	return node, ok
}

func (s *Service) GetConnectedNodes() map[string]*Node {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return maps.Clone(s.connectedNodes)
}

func (s *Service) SetConnectedNode(serverAddr string, node *Node) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.connectedNodes[serverAddr] = node
}

func (s *Service) DeleteConnectedNode(serverAddr string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.connectedNodes, serverAddr)
}

func NewService() *Service {
	return &Service{
		ShouldBroadcast: make(chan struct{}),
		connectedNodes:  make(map[string]*Node),
		pendingTxQ:      make(map[common.TxHashSized]*block_proto.Transaction),
	}
}
