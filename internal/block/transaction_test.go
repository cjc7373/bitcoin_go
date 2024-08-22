package block

import (
	"testing"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/utils"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	assert := assert.New(t)

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
	err := Sign(tx, w1.PrivateKey)
	tx.Id = hash(tx)
	assert.Nil(err)

	res, err := Verify(tx, unspentOutputs)
	assert.True(res)
	assert.Nil(err)

	// we tamper a signature, leaving the tx hash incorrect
	tx.VIn[0].Signature = tx.VIn[0].Signature[1:]
	res, err = Verify(tx, unspentOutputs)
	assert.False(res)
	assert.Equal(ErrInvalidHash, err)

	// we correct the hash
	tx.Id = hash(tx)
	res, err = Verify(tx, unspentOutputs)
	assert.False(res)
	assert.Equal(ErrInvalidSignature, err)
}
