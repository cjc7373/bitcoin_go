package network

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/cjc7373/bitcoin_go/internal/network/proto"
	"github.com/stretchr/testify/assert"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestDiscovery(t *testing.T) {
	serverAddr := ":12200"
	service := NewService()
	go func() {
		lis, _ := net.Listen("tcp", serverAddr)
		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		discovery := discoveryServer{s: service}
		proto.RegisterDiscoveryServer(grpcServer, &discovery)
		grpcServer.Serve(lis)
	}()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithNoProxy()}
	conn, err := grpc.Dial(serverAddr, opts...)
	assert.Nil(t, err)
	defer conn.Close()

	client := proto.NewDiscoveryClient(conn)
	_, err = client.RequestNodes(context.Background(), &proto.NodeRequest{Name: "foo"})
	fmt.Println(err)
	fmt.Println(service.connectedNodes)
	assert.Nil(t, err)
	assert.Len(t, service.connectedNodes, 1)
}
