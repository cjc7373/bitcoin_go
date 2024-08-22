package block

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/db"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
)

var _ = Describe("block test", func() {
	w1 := wallet.NewWallet()
	w2 := wallet.NewWallet()
	testConf.Wallets[testWalletName] = string(w1.EncodeToPEM())
	bdb := db.GetDB(&testConf)

	bc := NewBlockchain(bdb, w1.GetAddress())
	Reindex(bdb, bc)
	tx1, err := NewTransaction(bdb, w1, w2.GetAddress(), 100)
	Expect(err).NotTo(HaveOccurred())
	tx2, err := NewTransaction(bdb, w1, w2.GetAddress(), 200)
	Expect(err).NotTo(HaveOccurred())
	AddBlock(bdb, bc, []*block_proto.Transaction{tx1, tx2})
})
