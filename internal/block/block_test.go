package block

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
)

var _ = Describe("block test", func() {
	It("constructs block", func() {
		w1 := wallet.NewWallet()
		w2 := wallet.NewWallet()
		testConf.Wallets[testWalletName] = string(w1.EncodeToPEM())

		bc, err := NewBlockchain(testDB, w1.GetAddress())
		Expect(err).NotTo(HaveOccurred())
		Reindex(testDB, bc)
		tx1, err := NewTransaction(testDB, w1, w2.GetAddress(), 100)
		Expect(err).NotTo(HaveOccurred())
		tx2, err := NewTransaction(testDB, w1, w2.GetAddress(), 200)
		Expect(err).NotTo(HaveOccurred())
		AddBlock(testDB, bc, []*block_proto.Transaction{tx1, tx2})
	})
})
