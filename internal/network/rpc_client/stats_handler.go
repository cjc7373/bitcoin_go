package rpc_client

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"github.com/cjc7373/bitcoin_go/internal/network"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
)

type statsHandler struct {
	service *network.Service
	logger  *slog.Logger
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
			if err := DisconnectNode(h.service, addr); err != nil {
				h.logger.Error("disconnect error", "error", err)
			}
		}
	}
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
