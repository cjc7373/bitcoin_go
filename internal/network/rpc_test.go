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
	serverAddr1 := ":12200"
	service1 := NewService()
	done := make(chan error)
	go func() {
		service1.Serve(serverAddr1, done)
	}()

	serverAddr2 := ":12201"
	service2 := NewService()
	go func() {
		service2.Serve(serverAddr2, done)
	}()

	// wait server start
	time.Sleep(time.Microsecond * 100)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithNoProxy()}
	conn, err := grpc.Dial(serverAddr1, opts...)
	assert.Nil(t, err)

	client := proto.NewDiscoveryClient(conn)
	_, err = client.RequestNodes(context.Background(), &proto.Node{Name: "foo", Address: serverAddr2})
	fmt.Println(err)
	fmt.Println(service1.connectedNodes)
	assert.Nil(t, err)
	assert.Len(t, service1.connectedNodes, 1)

	conn.Close()
	// wait server handle conn close
	time.Sleep(time.Microsecond * 100)
	assert.Len(t, service1.connectedNodes, 0)
}
