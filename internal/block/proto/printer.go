package proto

import (
	"fmt"
)

type TransactionPretty Transaction

func (tx *TransactionPretty) String() string {
	rtn := fmt.Sprintf("Id: %x, ", tx.Id)
	rtn += "[]VIn{ "

	for _, input := range tx.VIn {
		rtn += fmt.Sprintf("{Txid: %x, Out: %d, Signature: %x, Pubkey: %x}, ", input.Txid, input.VoutIndex, input.Signature, input.PubKey)
	}
	rtn += "}, "

	rtn += "[]VOut{ "
	for _, output := range tx.VOut {
		rtn += fmt.Sprintf("{Value: %d, PubKeyHash: %x}, ", output.Value, output.PubKeyHash)
	}
	rtn += "} "
	return rtn
}

type UTXOSetPretty UTXOSet

func (utxoSet *UTXOSetPretty) String() string {
	rtn := "["
	for _, utxo := range utxoSet.UTXOs {
		rtn += fmt.Sprintf(" {txHash: %x, index: %v}, ", utxo.Transaction, utxo.OutputIndex)
	}
	rtn += "]"
	return rtn
}

type BlockchainPretty Blockchain

func (bc *BlockchainPretty) String() string {
	return fmt.Sprintf("Blockchain{tipHash: %x, height: %v}", bc.TipHash, bc.Height)
}
