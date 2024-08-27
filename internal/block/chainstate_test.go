package block

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
)

func newFakeBlockchain() *block_proto.Blockchain {
	tx1 := block_proto.Transaction{
		Id: []byte("1"),
		VIn: []*block_proto.TXInput{
			{VoutIndex: -1},
		},
		VOut: []*block_proto.TXOutput{
			{Value: 100},
			{Value: 200},
		},
	}
	tx2 := block_proto.Transaction{
		Id: []byte("2"),
		VIn: []*block_proto.TXInput{
			{Txid: []byte("1"), VoutIndex: 1},
		},
		VOut: []*block_proto.TXOutput{
			{Value: 1},
			{Value: 2},
			{Value: 3},
			{Value: 4},
		},
	}
	tx3 := block_proto.Transaction{
		Id: []byte("3"),
		VIn: []*block_proto.TXInput{
			{Txid: []byte("1"), VoutIndex: 0},
			{Txid: []byte("2"), VoutIndex: 0},
		},
		VOut: []*block_proto.TXOutput{
			{Value: 1},
			{Value: 2},
			{Value: 3},
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
			{Value: 1},
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
		bc := newFakeBlockchain()
		utxo := findUTXO(testDB, bc)
		// for k, v := range *utxo {
		// 	fmt.Println(k, v)
		// }

		Expect((*utxo)["4"][0].OriginalIndex).To(BeEquivalentTo(0))
		Expect((*utxo)["3"][0].OriginalIndex).To(BeEquivalentTo(0))
		Expect((*utxo)["2"][0].OriginalIndex).To(BeEquivalentTo(1))
		Expect((*utxo)["2"][1].OriginalIndex).To(BeEquivalentTo(3))
	})
})
