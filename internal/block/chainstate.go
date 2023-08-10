package block

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

const utxoBucket = "chainstate"

// in uxto bucket, we'll have:
// 32-byte tx hash -> []int, stores unspent output indexes in that tx

type UTXOSet struct {
	Blockchain *Blockchain
}

// return a map of UTXOs, key is tx id, value is a set of VoutIndex
func (u UTXOSet) FindUTXO() *map[string][]int {
	blockIter := u.Blockchain.Iterator()

	UTXO := make(map[string][]int)         // key is tx id, value is a set of VoutIndex
	spentOutputs := make(map[string][]int) // key is tx id, value is a set of VoutIndex

	for {
		block := blockIter.Next()
		if block == nil {
			break
		}

		for _, tx := range block.Transactions {
			id := string(tx.ID)

			for voutIndex := range tx.Vout {
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
					UTXO[id] = append(UTXO[id], voutIndex)
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

	utxo := u.FindUTXO()
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
