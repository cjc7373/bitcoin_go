package block

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
)

func newFakeBlockchain(addr wallet.Address) *block_proto.Blockchain {
	tx1 := block_proto.Transaction{
		Id: []byte("1"),
		VIn: []*block_proto.TXInput{
			{VoutIndex: -1},
		},
		VOut: []*block_proto.TXOutput{
			NewTXOutput(100, addr),
			NewTXOutput(200, addr),
		},
	}
	tx2 := block_proto.Transaction{
		Id: []byte("2"),
		VIn: []*block_proto.TXInput{
			{Txid: []byte("1"), VoutIndex: 1},
		},
		VOut: []*block_proto.TXOutput{
			NewTXOutput(1, addr),
			NewTXOutput(2, addr),
			NewTXOutput(3, addr),
			NewTXOutput(4, addr),
		},
	}
	tx3 := block_proto.Transaction{
		Id: []byte("3"),
		VIn: []*block_proto.TXInput{
			{Txid: []byte("1"), VoutIndex: 0},
			{Txid: []byte("2"), VoutIndex: 0},
		},
		VOut: []*block_proto.TXOutput{
			NewTXOutput(1, addr),
			NewTXOutput(2, addr),
			NewTXOutput(3, addr),
		},
	}
	tx4 := block_proto.Transaction{
		Id: []byte("4"),
		VIn: []*block_proto.TXInput{
			{Txid: []byte("3"), VoutIndex: 2},
			{Txid: []byte("3"), VoutIndex: 1},
			{Txid: []byte("2"), VoutIndex: 2},
		},
		VOut: []*block_proto.TXOutput{
			NewTXOutput(1, addr),
		},
	}
	bc := &block_proto.Blockchain{}
	var err error
	ctx := context.Background()
	// the UXTO will be (txid, outputIndex): (2, 1), (2, 3) (3, 0), (4, 0)
	_, err = AddBlock(ctx, testDB, bc, []*block_proto.Transaction{&tx1})
	Expect(err).NotTo(HaveOccurred())
	_, err = AddBlock(ctx, testDB, bc, []*block_proto.Transaction{&tx2})
	Expect(err).NotTo(HaveOccurred())
	_, err = AddBlock(ctx, testDB, bc, []*block_proto.Transaction{&tx3})
	Expect(err).NotTo(HaveOccurred())
	_, err = AddBlock(ctx, testDB, bc, []*block_proto.Transaction{&tx4})
	Expect(err).NotTo(HaveOccurred())

	return bc
}

var _ = Describe("chainstate test", func() {
	It("passes", func() {
		w1 := wallet.NewWallet()
		_ = newFakeBlockchain(w1.GetAddress())
		Expect(RebuildChainState(testDB)).To(Succeed())

		utxoSet, err := getUTXOSet(testDB, w1.GetAddress())
		Expect(err).NotTo(HaveOccurred())

		Expect(utxoSet.UTXOs[0].Transaction).To(BeEquivalentTo("4"))
		Expect(utxoSet.UTXOs[0].OutputIndex).To(BeEquivalentTo(0))

		Expect(utxoSet.UTXOs[0].Transaction).To(BeEquivalentTo("3"))
		Expect(utxoSet.UTXOs[0].OutputIndex).To(BeEquivalentTo(0))

		Expect(utxoSet.UTXOs[0].Transaction).To(BeEquivalentTo("2"))
		Expect(utxoSet.UTXOs[0].OutputIndex).To(BeEquivalentTo(3))

		Expect(utxoSet.UTXOs[0].Transaction).To(BeEquivalentTo("2"))
		Expect(utxoSet.UTXOs[0].OutputIndex).To(BeEquivalentTo(1))
	})
})
