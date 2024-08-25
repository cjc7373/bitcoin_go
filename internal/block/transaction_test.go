package block

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/utils"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
)

var _ = Describe("transaction test", func() {
	It("signs tx", func() {
		w1 := wallet.NewWallet()
		w2 := wallet.NewWallet()

		coinbaseTx := NewCoinbaseTransaction(w1.GetAddress(), nil)
		unspentOutputs := make(map[string][]TXOutputWithMetadata)
		unspentOutputs[string(coinbaseTx.Id)] = append(unspentOutputs[string(coinbaseTx.Id)], TXOutputWithMetadata{
			TXOutput:      coinbaseTx.VOut[0],
			OriginalIndex: 0,
		})

		in := block_proto.TXInput{Txid: coinbaseTx.Id, VoutIndex: 0}
		out := block_proto.TXOutput{Value: 100, PubKeyHash: utils.HashPubKey(w2.PublicKey)}

		tx := &block_proto.Transaction{
			Id:   nil,
			VIn:  []*block_proto.TXInput{&in},
			VOut: []*block_proto.TXOutput{&out},
		}
		Expect(Sign(tx, w1.PrivateKey)).To(Succeed())
		tx.Id = hashTx(tx)

		res, err := Verify(tx, unspentOutputs)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeTrue())

		// we tamper a signature, leaving the tx hash incorrect
		tx.VIn[0].Signature = tx.VIn[0].Signature[1:]
		res, err = Verify(tx, unspentOutputs)
		Expect(err).To(Equal(ErrInvalidHash))
		Expect(res).To(BeFalse())

		// we correct the hash
		tx.Id = hashTx(tx)
		res, err = Verify(tx, unspentOutputs)
		Expect(err).To(Equal(ErrInvalidSignature))
		Expect(res).To(BeFalse())
	})
})
