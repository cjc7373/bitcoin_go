package block

import (
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/common"
)

var blockchainKey = []byte("blockchain")

var ErrBlockchainNotExist = errors.New("blockchain not exists")
var ErrBlockchainAlreadyExist = errors.New("blockchain already exists")

func saveBlockchain(tx *bolt.Tx, bc *block_proto.Blockchain) error {
	b := tx.Bucket([]byte(common.BlockBucket))
	data, err := proto.Marshal(bc)
	if err != nil {
		return err
	}

	if err := b.Put(blockchainKey, data); err != nil {
		return err
	}
	return nil
}

func GetBlockchain(db *bolt.DB) (*block_proto.Blockchain, error) {
	bc := &block_proto.Blockchain{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(common.BlockBucket))
		v := b.Get(blockchainKey)
		if v == nil {
			return ErrBlockchainNotExist
		}
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

func NewBlockchain(bolt_db *bolt.DB, to string) (*block_proto.Blockchain, error) {
	if err := AddGenesisBlock(bolt_db, to); err != nil {
		return nil, err
	}

	return GetBlockchain(bolt_db)
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
