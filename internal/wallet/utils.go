package wallet

import (
	"fmt"

	"github.com/cjc7373/bitcoin_go/internal/utils"
)

func ReadOrCreateWalletFromConfig(conf *utils.Config) *Wallet {
	var w *Wallet
	if len(conf.Wallets) == 0 {
		w := NewWallet()
		fmt.Printf("Creating new wallet with address %v\n", w.GetAddress())
		conf.Wallets[DefatltWalletName] = string(w.EncodeToPEM())
		conf.Write()
	} else {
		w = NewWalletFromPEM([]byte(conf.Wallets[DefatltWalletName]))
		fmt.Printf("Using existing wallet with address %v\n", w.GetAddress())
	}
	return w
}
