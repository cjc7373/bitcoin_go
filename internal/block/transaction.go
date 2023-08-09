package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cjc7373/bitcoin_go/internal/utils"
)

type TXInput struct {
	Txid      []byte // ID of tx this input refers
	VoutIndex int    // index of an output in the tx
	Signature []byte
	PubKey    []byte
}

type TXOutput struct {
	// stores the number of satoshis, which is 0.00000001 BTC.
	// this is the smallest unit of currency in Bitcoin
	Value int64
	// we are not implementing the whole srcipt thing here, so just pubkey
	// pubkey hash is just pubkey hash, not an address
	PubKeyHash []byte
}

// NewTXOutput create a new TXOutput
// trim the address to only contain pubkey hash
func NewTXOutput(value int64, address string) *TXOutput {
	addressBytes := base58.Decode(address)
	pubkeyHash := addressBytes[1 : len(addressBytes)-4]

	txo := &TXOutput{Value: value, PubKeyHash: pubkeyHash}
	return txo
}

type Transaction struct {
	ID   []byte // hash of this tx
	Vin  []TXInput
	Vout []TXOutput
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	data, err := json.Marshal(&txCopy)
	if err != nil {
		panic(err)
	}
	hash = sha256.Sum256(data)

	return hash[:]
}

const subsidy = 10000

// create a new tx, which has an output to reward the miner
// this output
func NewCoinbaseTransaction(to string) *Transaction {
	// create an empty input to make the hash change every time
	randData := make([]byte, 10)
	_, err := rand.Read(randData)
	if err != nil {
		panic(err)
	}
	input := TXInput{nil, -1, nil, randData}
	output := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{input}, []TXOutput{*output}}
	tx.ID = tx.Hash()
	return &tx
}

// IsCoinbase checks whether the transaction is coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].VoutIndex == -1
}

// String returns a human-readable representation of a transaction
func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))

	for i, input := range tx.Vin {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.VoutIndex))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) error {
	if tx.IsCoinbase() {
		return nil
	}

	// bitcoin actually signs a trimmed copy of a tx, I don't know why
	// here I only sign an input
	for _, vin := range tx.Vin {
		prevTx := prevTXs[string(vin.Txid)]
		prevOutput := prevTx.Vout[vin.VoutIndex]
		pubkey := utils.EncodePubKey(&privKey)
		if !bytes.Equal(prevOutput.PubKeyHash, utils.HashPubKey(pubkey)) {
			return errors.New("pubkey not equal to previous output's pubkey hash")
		}

		vin.Signature = nil
		vin.PubKey = pubkey
		data, err := json.Marshal(&vin)
		if err != nil {
			panic(err)
		}
		hash := sha256.Sum256(data)
		sig, err := ecdsa.SignASN1(rand.Reader, &privKey, hash[:])
		if err != nil {
			panic(err)
		}
		vin.Signature = sig
	}
	return nil
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		vinCopy := vin
		vinCopy.Signature = nil

		prevOutput := prevTXs[string(vin.Txid)].Vout[vin.VoutIndex]
		if !bytes.Equal(prevOutput.PubKeyHash, utils.HashPubKey(vin.PubKey)) {
			return false
		}

		data, err := json.Marshal(&vinCopy)
		if err != nil {
			panic(err)
		}
		hash := sha256.Sum256(data)
		if !ecdsa.VerifyASN1(utils.ParsePubKey(vin.PubKey), hash[:], vin.Signature) {
			return false
		}
	}
	return true
}
