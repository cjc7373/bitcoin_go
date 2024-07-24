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
	"github.com/cjc7373/bitcoin_go/internal/wallet"
)

var (
	ErrPubKeyMismatch   = errors.New("pubkey not equal to previous output's pubkey hash")
	ErrInvalidSignature = errors.New("invalid input signature")
	ErrInvalidHash      = errors.New("invalid transaction hash")
)

type ErrNotEnoughFunds struct {
	need  int64
	found int64
}

func (err ErrNotEnoughFunds) Error() string {
	return fmt.Sprintf("the wallet do not have enough funds, need %v, found %v", err.need, err.found)
}

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

// hash returns the hash of the Transaction
func (tx *Transaction) hash() []byte {
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
func NewCoinbaseTransaction(to string, data []byte) *Transaction {
	if data != nil {
		// create an empty input to make the hash change every time
		data = make([]byte, 10)
		_, err := rand.Read(data)
		if err != nil {
			panic(err)
		}
	}

	input := TXInput{nil, -1, nil, data}
	output := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{input}, []TXOutput{*output}}
	tx.ID = tx.hash()
	return &tx
}

func NewTransaction(w *wallet.Wallet, to string, amount int64, uxtoSet *UTXOSet) (*Transaction, error) {
	unspentOutputs, foundAmount := uxtoSet.FindSpendableOutputs(utils.HashPubKey(w.PublicKey), amount)
	if foundAmount < amount {
		return nil, ErrNotEnoughFunds{need: amount, found: foundAmount}
	}

	var inputs []TXInput
	for txID, outputs := range unspentOutputs {
		for _, output := range outputs {
			inputs = append(inputs, TXInput{
				Txid:      []byte(txID),
				VoutIndex: output.OriginalIndex,
			})
		}
	}
	outputs := []TXOutput{{amount, nil}}
	// take the change
	if foundAmount > amount {
		outputs = append(outputs, TXOutput{foundAmount - amount, utils.HashPubKey(w.PublicKey)})
	}
	tx := Transaction{nil, inputs, outputs}
	if err := tx.Sign(w.PrivateKey); err != nil {
		return nil, err
	}
	tx.ID = tx.hash()
	return &tx, nil
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
		lines = append(lines, fmt.Sprintf("       PubKeyHash: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

// in Sign() function we do not need to verify the pubkey of vin
// because the transaction will always be valid
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) error {
	if tx.IsCoinbase() {
		return nil
	}

	// bitcoin actually signs a trimmed copy of a tx, I don't know why
	// here I only sign an input
	for index := range tx.Vin {
		pubkey := utils.EncodePubKey(&privKey)

		tx.Vin[index].Signature = nil
		tx.Vin[index].PubKey = pubkey
		data, err := json.Marshal(&tx.Vin[index])
		if err != nil {
			return err
		}
		hash := sha256.Sum256(data)
		sig, err := ecdsa.SignASN1(rand.Reader, &privKey, hash[:])
		if err != nil {
			return err
		}
		tx.Vin[index].Signature = sig
	}
	return nil
}

func findOutputInUXTO(unspentOutputs *map[string][]TXOutputWithMetadata, txid string, outputIndex int) *TXOutputWithMetadata {
	for id, outputs := range *unspentOutputs {
		for _, output := range outputs {
			if txid == id && output.OriginalIndex == outputIndex {
				return &output
			}
		}
	}
	return nil
}

func (tx *Transaction) Verify(prevOutputs map[string][]TXOutputWithMetadata) (bool, error) {
	if tx.IsCoinbase() {
		return true, nil
	}

	txHash := tx.hash()
	if !bytes.Equal(tx.ID, txHash) {
		return false, ErrInvalidHash
	}

	for _, vinCopy := range tx.Vin {
		sig := vinCopy.Signature
		vinCopy.Signature = nil

		prevOutput := findOutputInUXTO(&prevOutputs, string(vinCopy.Txid), vinCopy.VoutIndex)
		if !bytes.Equal(prevOutput.PubKeyHash, utils.HashPubKey(vinCopy.PubKey)) {
			return false, ErrPubKeyMismatch
		}

		data, err := json.Marshal(&vinCopy)
		if err != nil {
			panic(err)
		}
		hash := sha256.Sum256(data)
		if !ecdsa.VerifyASN1(utils.ParsePubKey(vinCopy.PubKey), hash[:], sig) {
			return false, ErrInvalidSignature
		}
	}
	return true, nil
}
