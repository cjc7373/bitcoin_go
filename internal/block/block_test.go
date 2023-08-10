package block

import (
	"os"
	"testing"

	"github.com/cjc7373/bitcoin_go/internal/utils"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	assert := assert.New(t)

	dbPath := "blockchain_test.db"
	conf := utils.Config{
		DBPath:  dbPath,
		Wallets: map[string]string{},
	}
	t.Cleanup(func() {
		os.Remove(dbPath)
	})

	var defatltWalletName = "default"
	w1 := wallet.NewWallet()
	w2 := wallet.NewWallet()
	conf.Wallets[defatltWalletName] = string(w1.EncodeToPEM())

	bc := NewBlockchain(&conf, w1.GetAddress())
	utxoSet := UTXOSet{bc}
	utxoSet.Reindex()
	tx1, err := NewTransaction(w1, w2.GetAddress(), 100, &utxoSet)
	assert.Nil(err)
	tx2, err := NewTransaction(w1, w2.GetAddress(), 200, &utxoSet)
	assert.Nil(err)
	bc.AddBlock(&[]Transaction{*tx1, *tx2})
}
