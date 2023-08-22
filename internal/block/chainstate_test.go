package block

import (
	"fmt"
	"testing"

	"github.com/cjc7373/bitcoin_go/internal/utils"
	"github.com/stretchr/testify/assert"
)

type dumbIterator struct {
	cur    int
	blocks []*Block
}

func (it *dumbIterator) Next() bool {
	if it.cur < 0 {
		return false
	}
	it.cur -= 1
	return true
}

func (it *dumbIterator) Elem() *Block {
	return it.blocks[it.cur+1]
}

func newFakeBlockchain() *Blockchain {
	tx1 := Transaction{
		ID: []byte("1"),
		Vin: []TXInput{
			{nil, -1, nil, nil},
		},
		Vout: []TXOutput{
			{100, nil},
			{200, nil},
		},
	}
	tx2 := Transaction{
		ID: []byte("2"),
		Vin: []TXInput{
			{[]byte("1"), 1, nil, nil},
		},
		Vout: []TXOutput{
			{1, nil},
			{2, nil},
			{3, nil},
			{4, nil},
		},
	}
	tx3 := Transaction{
		ID: []byte("3"),
		Vin: []TXInput{
			{[]byte("1"), 0, nil, nil},
			{[]byte("2"), 0, nil, nil},
		},
		Vout: []TXOutput{
			{1, nil},
			{2, nil},
			{3, nil},
		},
	}
	tx4 := Transaction{
		ID: []byte("4"),
		Vin: []TXInput{
			{[]byte("3"), 2, nil, nil},
			{[]byte("3"), 1, nil, nil},
			{[]byte("2"), 2, nil, nil},
		},
		Vout: []TXOutput{
			{1, nil},
		},
	}
	// the UXTO will be (txid, outputIndex): (2, 1), (2, 3) (3, 0), (4, 0)
	b1 := NewBlock(&[]Transaction{tx1}, []byte{})
	b2 := NewBlock(&[]Transaction{tx2}, b1.Hash)
	b3 := NewBlock(&[]Transaction{tx3}, b2.Hash)
	b4 := NewBlock(&[]Transaction{tx4}, b3.Hash)

	bc := Blockchain{
		TipHash: b4.Hash,
		Height:  4,
		DB:      nil,
		NewBlockIterator: func() utils.Iterator[*Block] {
			return &dumbIterator{cur: 3, blocks: []*Block{b1, b2, b3, b4}}
		},
	}
	return &bc
}

func TestFindUTXO(t *testing.T) {
	assert := assert.New(t)

	bc := newFakeBlockchain()
	utxoSet := UTXOSet{Blockchain: bc}
	utxo := utxoSet.findUTXO()
	for k, v := range *utxo {
		fmt.Println(k, v)
	}
	assert.Equal(0, (*utxo)["4"][0].OriginalIndex)
	assert.Equal(0, (*utxo)["3"][0].OriginalIndex)
	assert.Equal(1, (*utxo)["2"][0].OriginalIndex)
	assert.Equal(3, (*utxo)["2"][1].OriginalIndex)

}
