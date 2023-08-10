package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/cjc7373/bitcoin_go/internal/block"
	"github.com/cjc7373/bitcoin_go/internal/utils"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
)

func main() {
	configPath := flag.String("config", "config.yaml", "config path")
	flag.Parse()
	configPathAbs, _ := filepath.Abs(*configPath)
	fmt.Printf("Using config: %v\n", configPathAbs)
	conf := utils.ParseConfig(*configPath)

	var w wallet.Wallet
	var defatltWalletName = "default"
	if len(conf.Wallets) == 0 {
		w := wallet.NewWallet()
		fmt.Printf("Creating new wallet with address %v\n", w.GetAddress())
		conf.Wallets[defatltWalletName] = string(w.EncodeToPEM())
		conf.WriteToFile(*configPath)
	} else {
		w = *wallet.NewWalletFromPEM([]byte(conf.Wallets[defatltWalletName]))
		fmt.Printf("Using existing wallet with address %v\n", w.GetAddress())
	}

	bc := block.NewBlockchain(conf, w.GetAddress())

	// bc.AddBlock("Send 1 BTC to Ivan")
	// bc.AddBlock("Send 2 more BTC to Ivan")

	bc.PrintChain()
}
