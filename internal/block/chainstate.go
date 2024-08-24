package block

import (
	"bytes"
	"encoding/json"

	bolt "go.etcd.io/bbolt"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/common"
)

// in uxto bucket, we'll have:
// 32-byte tx hash -> []int, stores unspent output indexes in that tx

type TXOutputWithMetadata struct {
	*block_proto.TXOutput
	// there are only unspent outputs in UTXO
	// so we need this field to identify its original VoutIndex
	OriginalIndex int32
}

// find someone's enough outputs to make the tx
// FIXME: this function needs to iterate chainstate bucket, could be slow
// return a map of someone's UXTOs, key is tx id, value is a set of TXOutputs
func FindSpendableOutputs(db *bolt.DB, pubkeyHash []byte, amount int64) (unspentOutputs map[string][]TXOutputWithMetadata, found int64) {
	unspentOutputs = make(map[string][]TXOutputWithMetadata)
	var accumulated int64 = 0

	err := db.View(func(blotTx *bolt.Tx) error {
		b := blotTx.Bucket([]byte(common.UTXOBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outputs := new([]TXOutputWithMetadata)
			err := json.Unmarshal(v, outputs)
			if err != nil {
				return err
			}
			for _, output := range *outputs {
				if bytes.Equal(output.PubKeyHash, pubkeyHash) && accumulated < amount {
					accumulated += output.Value
					unspentOutputs[string(k)] = append(unspentOutputs[string(k)], output)
				}
			}
			if accumulated > amount {
				break
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return unspentOutputs, accumulated
}

// return a map of UTXOs, key is tx id, value is a set of TXOutputs
func findUTXO(db *bolt.DB, bc *block_proto.Blockchain) *map[string][]TXOutputWithMetadata {
	UTXO := make(map[string][]TXOutputWithMetadata) // key is tx id, value is a set of VoutIndex
	spentOutputs := make(map[string][]int32)        // key is tx id, value is a set of VoutIndex

	for block := range AllBlocks(db, bc.TipHash) {
		for _, tx := range block.Transactions {
			id := string(tx.Id)

			for voutIndex, output := range tx.VOut {
				spent := false
				// if spentOutputs[id] doesn't exist, this output can't be spent
				// because we iterate the chain from end to start
				// Note: here we are assuming that an output will not be spent in the same block
				if _, exist := spentOutputs[id]; exist {
					// iterator spentOutputs[id] to check if this output is spent
					for _, spentOutput := range spentOutputs[id] {
						if spentOutput == int32(voutIndex) {
							spent = true
							break
						}
					}
				}

				if !spent {
					UTXO[id] = append(UTXO[id], TXOutputWithMetadata{TXOutput: output, OriginalIndex: int32(voutIndex)})
				}
			}

			// coinbase tx doesn't spend any outputs
			if !IsCoinbase(tx) {
				for _, input := range tx.VIn {
					spentOutputs[string(input.Txid)] = append(spentOutputs[string(input.Txid)], input.VoutIndex)
				}
			}
		}
	}

	return &UTXO
}

// rebuild UXTO set
func Reindex(db *bolt.DB, bc *block_proto.Blockchain) {
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(common.UTXOBucket))
		if err != nil && err != bolt.ErrBucketNotFound {
			panic(err)
		}

		_, err = tx.CreateBucket([]byte(common.UTXOBucket))
		if err != nil {
			panic(err)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	utxo := findUTXO(db, bc)
	// update utxo in db
	err = db.Update(func(dbTx *bolt.Tx) error {
		b := dbTx.Bucket([]byte(common.UTXOBucket))

		for k, v := range *utxo {
			data, err := json.Marshal(&v)
			if err != nil {
				return err
			}
			err = b.Put([]byte(k), data)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
