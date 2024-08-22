package block

import (
	"encoding/json"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
)

var blockchainBucket = []byte("blockchain")
var blockchainKey = []byte("blockchain")

func GetBlockchain(db *bolt.DB) (*block_proto.Blockchain, error) {
	bc := &block_proto.Blockchain{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(blockchainBucket)
		v := b.Get(blockchainKey)
		if err := proto.Unmarshal(v, bc); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bc, nil
}

func NewBlockchain(bolt_db *bolt.DB, to string) *block_proto.Blockchain {
	var tip []byte

	err := bolt_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))

		if b == nil {
			genesis := NewGenesisBlock(to)
			log.Printf("Created genesis block with hash %x\n", genesis.Hash)
			b, err := tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				panic(err)
			}

			data, err := json.Marshal(genesis)
			if err != nil {
				panic(err)
			}

			err = b.Put(genesis.Hash, data)
			if err != nil {
				panic(err)
			}

			err = b.Put([]byte(lastBlockKey), genesis.Hash)
			if err != nil {
				panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte(lastBlockKey))
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return &block_proto.Blockchain{TipHash: tip, Height: 1}
}

func PrintChain(db *bolt.DB, bc *block_proto.Blockchain) {
	fmt.Println("Printing chain...")

	for block := range AllBlocks(db, bc.TipHash) {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Transactions: \n")
		for _, tx := range block.Transactions {
			fmt.Println(&tx)
		}
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}

}
