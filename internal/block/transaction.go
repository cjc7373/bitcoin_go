package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
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

// NewTXOutput create a new TXOutput
// trim the address to only contain pubkey hash
func NewTXOutput(value int64, address wallet.Address) *block_proto.TXOutput {
	pubkeyHash := wallet.GetPubKey(address)

	txo := &block_proto.TXOutput{Value: value, PubKeyHash: pubkeyHash}
	return txo
}

// hashTx returns the hashTx of the Transaction
func hashTx(tx *block_proto.Transaction) []byte {
	var hash [32]byte

	// hashTx will be used in both partial tx and finished tx,
	// so tx's hash might already exists. Remove it first.
	txCopy := proto.Clone(tx).(*block_proto.Transaction)
	txCopy.Id = []byte{}

	data, err := proto.Marshal(txCopy)
	if err != nil {
		panic(err)
	}
	hash = sha256.Sum256(data)

	return hash[:]
}

// reward for the miner
const subsidy = 10000

// create a new tx, which has an output to reward the miner
// this output
func NewCoinbaseTransaction(to wallet.Address, data []byte) *block_proto.Transaction {
	if data == nil {
		// create an empty input to make the hash change every time
		data = make([]byte, 10)
		_, err := rand.Read(data)
		if err != nil {
			panic(err)
		}
	}

	input := block_proto.TXInput{
		Txid:      nil,
		VoutIndex: -1,
		Signature: nil,
		PubKey:    data,
	}
	output := NewTXOutput(subsidy, to)
	tx := block_proto.Transaction{
		Id:   nil,
		VIn:  []*block_proto.TXInput{&input},
		VOut: []*block_proto.TXOutput{output},
	}
	tx.Id = hashTx(&tx)
	return &tx
}

func NewTransaction(db *bolt.DB, w *wallet.Wallet, to wallet.Address, amount int64) (*block_proto.Transaction, error) {
	unspentOutputs, foundAmount := FindSpendableOutputs(db, utils.HashPubKey(w.PublicKey), amount)
	if foundAmount < amount {
		return nil, ErrNotEnoughFunds{need: amount, found: foundAmount}
	}

	var inputs []*block_proto.TXInput
	for txID, outputs := range unspentOutputs {
		for _, output := range outputs {
			inputs = append(inputs, &block_proto.TXInput{
				Txid:      []byte(txID),
				VoutIndex: output.OriginalIndex,
			})
		}
	}
	outputs := []*block_proto.TXOutput{NewTXOutput(amount, to)}
	// take the change
	if foundAmount > amount {
		outputs = append(outputs, &block_proto.TXOutput{
			Value:      foundAmount - amount,
			PubKeyHash: utils.HashPubKey(w.PublicKey),
		})
	}
	tx := &block_proto.Transaction{
		Id:   nil,
		VIn:  inputs,
		VOut: outputs,
	}
	if err := Sign(tx, w.PrivateKey); err != nil {
		return nil, err
	}
	tx.Id = hashTx(tx)
	return tx, nil
}

// IsCoinbase checks whether the transaction is coinbase
func IsCoinbase(tx *block_proto.Transaction) bool {
	return len(tx.VIn) == 1 && len(tx.VIn[0].Txid) == 0 && tx.VIn[0].VoutIndex == -1
}

// in Sign() function we do not need to verify the pubkey of vin
// because the transaction will always be valid
func Sign(tx *block_proto.Transaction, privKey ecdsa.PrivateKey) error {
	if IsCoinbase(tx) {
		return nil
	}

	// bitcoin actually signs a trimmed copy of a tx, I don't know why
	// here I only sign an input
	for index := range tx.VIn {
		pubkey := utils.EncodePubKey(&privKey)

		tx.VIn[index].Signature = nil
		tx.VIn[index].PubKey = pubkey
		data, err := proto.Marshal(tx.VIn[index])
		if err != nil {
			return err
		}
		hash := sha256.Sum256(data)
		sig, err := ecdsa.SignASN1(rand.Reader, &privKey, hash[:])
		if err != nil {
			return err
		}
		tx.VIn[index].Signature = sig
	}
	return nil
}

func findOutputInUXTO(unspentOutputs *map[string][]TXOutputWithMetadata, txid string, outputIndex int32) *TXOutputWithMetadata {
	for id, outputs := range *unspentOutputs {
		for _, output := range outputs {
			if txid == id && output.OriginalIndex == outputIndex {
				return &output
			}
		}
	}
	return nil
}

func Verify(tx *block_proto.Transaction, prevOutputs map[string][]TXOutputWithMetadata) (bool, error) {
	if IsCoinbase(tx) {
		return true, nil
	}

	txHash := hashTx(tx)
	if !bytes.Equal(tx.Id, txHash) {
		return false, ErrInvalidHash
	}

	for _, vin := range tx.VIn {
		vinCopy := proto.Clone(vin).(*block_proto.TXInput)
		sig := vinCopy.Signature
		vinCopy.Signature = nil

		prevOutput := findOutputInUXTO(&prevOutputs, string(vinCopy.Txid), vinCopy.VoutIndex)
		if !bytes.Equal(prevOutput.PubKeyHash, utils.HashPubKey(vinCopy.PubKey)) {
			return false, ErrPubKeyMismatch
		}

		data, err := proto.Marshal(vinCopy)
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
