package wallet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAddress(t *testing.T) {
	wallet := NewWallet()
	address := wallet.GetAddress()
	fmt.Println(address)
	// https://bitcoin.stackexchange.com/questions/36944/what-are-the-minimum-and-maximum-lengths-of-a-mainnet-bitcoin-address
	assert.GreaterOrEqual(t, len(address), 26)
}

func TestEncodePrivateKey(t *testing.T) {
	assert := assert.New(t)
	wallet := NewWallet()

	pemEncoded := wallet.EncodeToPEM()
	w2 := NewWalletFromPEM(pemEncoded)

	assert.Equal(wallet, w2)
}
