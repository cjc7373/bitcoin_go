package network

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cjc7373/bitcoin_go/internal/network/proto"
	"github.com/stretchr/testify/assert"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestDiscovery(t *testing.T) {
	serverAddr := ":12200"
	service := NewService()
	done := make(chan error)
	go func() {
		service.Serve(serverAddr, done)
	}()

	// wait server start
	time.Sleep(time.Microsecond * 100)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithNoProxy()}
	conn, err := grpc.Dial(serverAddr, opts...)
	assert.Nil(t, err)

	client := proto.NewDiscoveryClient(conn)
	_, err = client.RequestNodes(context.Background(), &proto.NodeRequest{Name: "foo"})
	fmt.Println(err)
	fmt.Println(service.connectedNodes)
	assert.Nil(t, err)
	assert.Len(t, service.connectedNodes, 1)

	conn.Close()
	// wait server handle conn close
	time.Sleep(time.Microsecond * 100)
	assert.Len(t, service.connectedNodes, 0)
}
