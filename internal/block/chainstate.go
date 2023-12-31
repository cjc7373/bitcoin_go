package block

import (
	"bytes"
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

const utxoBucket = "chainstate"

// in uxto bucket, we'll have:
// 32-byte tx hash -> []int, stores unspent output indexes in that tx

type UTXOSet struct {
	Blockchain *Blockchain
}

type TXOutputWithMetadata struct {
	TXOutput
	// there are only unspent outputs in UTXO
	// so we need this field to identify its original VoutIndex
	OriginalIndex int
}

// find someone's enough outputs to make the tx
// FIXME: this function needs to iterate chainstate bucket, could be slow
// return a map of someone's UXTOs, key is tx id, value is a set of TXOutputs
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int64) (unspentOutputs map[string][]TXOutputWithMetadata, found int64) {
	unspentOutputs = make(map[string][]TXOutputWithMetadata)
	var accumulated int64 = 0

	err := u.Blockchain.DB.View(func(blotTx *bolt.Tx) error {
		b := blotTx.Bucket([]byte(utxoBucket))
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
func (u UTXOSet) findUTXO() *map[string][]TXOutputWithMetadata {
	blockIter := u.Blockchain.NewBlockIterator()

	UTXO := make(map[string][]TXOutputWithMetadata) // key is tx id, value is a set of VoutIndex
	spentOutputs := make(map[string][]int)          // key is tx id, value is a set of VoutIndex

	for blockIter.Next() {
		block := blockIter.Elem()

		for _, tx := range block.Transactions {
			id := string(tx.ID)

			for voutIndex, output := range tx.Vout {
				spent := false
				// if spentOutputs[id] doesn't exist, this output can't be spent
				// because we iterate the chain from end to start
				// Note: here we are assuming that an output will not be spent in the same block
				if _, exist := spentOutputs[id]; exist {
					// iterator spentOutputs[id] to check if this output is spent
					for _, spentOutput := range spentOutputs[id] {
						if spentOutput == voutIndex {
							spent = true
							break
						}
					}
				}

				if !spent {
					UTXO[id] = append(UTXO[id], TXOutputWithMetadata{TXOutput: output, OriginalIndex: voutIndex})
				}
			}

			// coinbase tx doesn't spend any outputs
			if !tx.IsCoinbase() {
				for _, input := range tx.Vin {
					spentOutputs[string(input.Txid)] = append(spentOutputs[string(input.Txid)], input.VoutIndex)
				}
			}
		}
	}

	return &UTXO
}

// rebuild UXTO set
func (u UTXOSet) Reindex() {
	db := u.Blockchain.DB

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(utxoBucket))
		if err != nil && err != bolt.ErrBucketNotFound {
			panic(err)
		}

		_, err = tx.CreateBucket([]byte(utxoBucket))
		if err != nil {
			panic(err)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	utxo := u.findUTXO()
	// update utxo in db
	err = db.Update(func(dbTx *bolt.Tx) error {
		b := dbTx.Bucket([]byte(utxoBucket))

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
