package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/cjc7373/bitcoin_go/internal/block"
	"github.com/cjc7373/bitcoin_go/internal/utils"
)

func main() {
	configPath := flag.String("config", "config.yaml", "config path")
	flag.Parse()
	configPathAbs, _ := filepath.Abs(*configPath)
	fmt.Printf("Using config: %v\n", configPathAbs)
	conf := utils.ParseConfig(*configPath)

	bc := block.NewBlockchain(conf)

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
