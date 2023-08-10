package block

import (
	"testing"

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
	unspentOutputs[string(coinbaseTx.ID)] = append(unspentOutputs[string(coinbaseTx.ID)], TXOutputWithMetadata{
		TXOutput:      coinbaseTx.Vout[0],
		OriginalIndex: 0,
	})

	in := TXInput{Txid: coinbaseTx.ID, VoutIndex: 0}
	out := TXOutput{Value: 100, PubKeyHash: utils.HashPubKey(w2.PublicKey)}

	tx := Transaction{nil, []TXInput{in}, []TXOutput{out}}
	err := tx.Sign(w1.PrivateKey)
	tx.ID = tx.hash()
	assert.Nil(err)

	res, err := tx.Verify(unspentOutputs)
	assert.True(res)
	assert.Nil(err)

	// we tamper a signature, leaving the tx hash incorrect
	tx.Vin[0].Signature = tx.Vin[0].Signature[1:]
	res, err = tx.Verify(unspentOutputs)
	assert.False(res)
	assert.Equal(ErrInvalidHash, err)

	// we correct the hash
	tx.ID = tx.hash()
	res, err = tx.Verify(unspentOutputs)
	assert.False(res)
	assert.Equal(ErrInvalidSignature, err)
}
