package rpc_server

import (
	"errors"
	"log/slog"
	"net"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/db"
	"github.com/cjc7373/bitcoin_go/internal/network"
	"github.com/cjc7373/bitcoin_go/internal/network/proto"
	"github.com/cjc7373/bitcoin_go/internal/network/rpc_client"
	"github.com/cjc7373/bitcoin_go/internal/utils"
)

type BitcoinServer struct {
	RPCClient rpc_client.RPCClient
	s         *network.Service
	logger    *slog.Logger
	config    *utils.Config

	// whether server has exited
	done       chan error
	grpcServer *grpc.Server

	DB *bolt.DB

	// persistant data
	blockchain block_proto.Blockchain

	proto.UnimplementedBitcoinServer
}

func NewRPCServer(logger *slog.Logger, config *utils.Config) BitcoinServer {
	service := network.NewService()
	db := db.OpenDB(config)

	return BitcoinServer{
		s:         service,
		logger:    logger,
		RPCClient: rpc_client.NewRPCClient(service, logger, db),
		config:    config,

		done: make(chan error),

		DB: db,
	}
}

func (d *BitcoinServer) Serve() error {
	lis, err := net.Listen("tcp", d.config.ListenAddr)
	if err != nil {
		return errors.Join(errors.New("listen failed"), err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	d.grpcServer = grpcServer
	proto.RegisterBitcoinServer(grpcServer, d)
	go d.RPCClient.HandleBroadcast()
	go func() {
		d.logger.Info("serving", "addr", d.config.ListenAddr)
		d.done <- grpcServer.Serve(lis)
	}()
	return nil
}

func (d *BitcoinServer) Stop() {
	d.grpcServer.GracefulStop()
	d.DB.Close()
}

func (d *BitcoinServer) GetConnectedNodes() map[string]*network.Node {
	return d.s.GetConnectedNodes()
}

func (d *BitcoinServer) DisconnectNode(address string) error {
	return rpc_client.DisconnectNode(d.s, address)
}
