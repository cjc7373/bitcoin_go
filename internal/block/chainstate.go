package block

import (
	"bytes"
	"encoding/json"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/common"
	bitcoin_db "github.com/cjc7373/bitcoin_go/internal/db"
	"github.com/cjc7373/bitcoin_go/internal/utils"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
)

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

func AddBlockToChainstate() {
	//TODO
}

// if multiple blocks will be reverted, they should be in the reverse order
// of the blockchain
func RevertBlockFromChainstate() {
	//TODO
}

// Normally this should not be called
func RebuildChainState(db *bolt.DB) error {
	bc, err := GetBlockchain(db)
	if err != nil {
		return err
	}
	if err := bitcoin_db.DeleteCacheBucket(db); err != nil {
		return err
	}

	// map of every address's utxo set
	allUTXOsets := make(map[utils.PubKeyHashSized]*block_proto.UTXOSet)
	spentOutputs := make(map[string]utils.Set[int]) // key is tx id, value is a set of VoutIndex

	for block := range AllBlocks(db, bc.TipHash) {
		for _, tx := range block.Transactions {
			if err := db.Update(func(dbTx *bolt.Tx) error {
				b := dbTx.Bucket(common.TransactionBucket)
				v, err := proto.Marshal(tx)
				if err != nil {
					return err
				}
				return b.Put(tx.Id, v)
			}); err != nil {
				return err
			}

			id := string(tx.Id)

			for voutIndex, output := range tx.VOut {
				spent := false
				// if spentOutputs[id] doesn't exist, this output can't be spent
				// because we iterate the chain from end to start
				// Note: here we are assuming that an output will not be spent in the same block
				if _, exist := spentOutputs[id]; exist {
					// iterator spentOutputs[id] to check if this output is spent
					if spentOutputs[id].Has(voutIndex) {
						spent = true
					}
				}

				if !spent {
					h := utils.PubKeyHashSized(output.PubKeyHash)
					if allUTXOsets[h] == nil {
						allUTXOsets[h] = &block_proto.UTXOSet{
							UTXOs: make([]*block_proto.UTXO, 0),
						}
					}
					allUTXOsets[h].UTXOs = append(allUTXOsets[h].UTXOs, &block_proto.UTXO{
						Transaction: tx.Id,
						OutputIndex: int32(voutIndex),
					})
				}
			}

			// coinbase tx doesn't spend any outputs
			if !IsCoinbase(tx) {
				for _, input := range tx.VIn {
					spentOutputs[string(input.Txid)].Insert(int(input.VoutIndex))
				}
			}
		}
	}

	return db.Update(func(dbTx *bolt.Tx) error {
		b := dbTx.Bucket(common.UTXOBucket)
		for pubKey, utxoSet := range allUTXOsets {
			v, err := proto.Marshal(utxoSet)
			if err != nil {
				return err
			}
			if err := b.Put(pubKey[:], v); err != nil {
				return err
			}
		}
		return nil
	})
}

func GetTransaction(db *bolt.DB, hash []byte) (*block_proto.Transaction, error) {
	tx := &block_proto.Transaction{}
	if err := db.View(func(dbTx *bolt.Tx) error {
		b := dbTx.Bucket(common.TransactionBucket)
		v := b.Get(hash)
		if v == nil {
			return nil
		}
		return proto.Unmarshal(v, tx)
	}); err != nil {
		return nil, err
	}
	return tx, nil
}

func getUTXOSet(db *bolt.DB, addr wallet.Address) (*block_proto.UTXOSet, error) {
	utxoSet := &block_proto.UTXOSet{}
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(common.UTXOBucket)
		v := b.Get(wallet.GetPubKey(addr))
		if v == nil {
			return nil
		}
		return proto.Unmarshal(v, utxoSet)
	}); err != nil {
		return nil, err
	}
	return utxoSet, nil
}

func GetBalance(db *bolt.DB, addr wallet.Address) (amount int64, err error) {
	utxoSet, err := getUTXOSet(db, addr)
	if err != nil {
		return
	}

	for _, utxo := range utxoSet.UTXOs {
		tx, err := GetTransaction(db, utxo.Transaction)
		if err != nil {
			return 0, err
		}
		amount += tx.VOut[utxo.OutputIndex].Value
	}
	return
}
