package rpc_server

import (
	"context"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/network/proto"
)

func (d *BitcoinServer) SendNodes(ctx context.Context, nodes *proto.Nodes) (*proto.Empty, error) {
	connectedNodesUpdated := false
	// TODO: use a channel to de-couple RPCClient.ConnectNode
	for _, node := range nodes.Nodes {
		connected, err := d.RPCClient.ConnectNode(node.Address, node.Name, d.config.ListenAddr, d.config.NodeName)
		if err != nil {
			return nil, err
		}
		if connected {
			connectedNodesUpdated = true
		}
	}

	if connectedNodesUpdated {
		// don't block on this channel
		select {
		case d.s.ShouldBroadcast <- struct{}{}:
		default:
			d.logger.Info("ShouldBroadcast not ready to receive, skipping")
		}
	}

	return &proto.Empty{}, nil
}

func (d *BitcoinServer) SendChainMetadata(ctx context.Context, bc *block_proto.Blockchain) (*proto.Empty, error) {
	return &proto.Empty{}, nil
}
