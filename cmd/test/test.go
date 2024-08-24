package main

import (
	"fmt"
	"time"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
)

func main() {
	b := &block_proto.Block{
		Timestamp: time.Now().Unix(),
		Hash:      []byte("1231232"),
	}
	fmt.Println(b)
}
