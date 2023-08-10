package block

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/cjc7373/bitcoin_go/internal/db"
	"github.com/cjc7373/bitcoin_go/internal/utils"
)

const blockBucket = "block"
const lastBlock = "last_block"

// in block bucket, we'll have:
// 32-byte block hash -> block data, encoded by json
// "last_block" -> the hash of the last block in a chain

type Block struct {
	Timestamp     int64 // a Unix timestamp
	Transactions  []Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(txs *[]Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), *txs, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

type Blockchain struct {
	TipHash  []byte // top block hash
	Height   int64
	DB       *bolt.DB
	Iterator func() BlockchainIterator
}

func (bc *Blockchain) AddBlock(txs *[]Transaction) {
	newBlock := NewBlock(txs, bc.TipHash)
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		data, err := json.Marshal(newBlock)
		if err != nil {
			panic(err)
		}

		err = b.Put(newBlock.Hash, data)
		if err != nil {
			panic(err)
		}

		err = b.Put([]byte(lastBlock), newBlock.Hash)
		if err != nil {
			panic(err)
		}

		bc.TipHash = newBlock.Hash

		return nil
	})

	if err != nil {
		panic(err)
	}
}

// FIXME: this interface is mainly for the ease of testing
// I wonder if there's a better way
type BlockchainIterator interface {
	Next() *Block
}

type BlockchainIteratorImpl struct {
	curHash []byte
	db      *bolt.DB
}

func (bci *BlockchainIteratorImpl) Next() *Block {
	if len(bci.curHash) == 0 {
		return nil
	}

	var block Block
	err := bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))

		data := b.Get(bci.curHash)
		err := json.Unmarshal(data, &block)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	bci.curHash = block.PrevBlockHash
	return &block
}

func (bc *Blockchain) PrintChain() {
	fmt.Println("Printing chain...")
	iter := bc.Iterator()

	for {
		block := iter.Next()

		if block == nil {
			break
		}

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Transactions: \n")
		for _, tx := range block.Transactions {
			fmt.Println(&tx)
		}
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}

}

func NewGenesisBlock(to string) *Block {
	return NewBlock(&[]Transaction{*NewCoinbaseTransaction(to, []byte("Genesis Block"))}, nil)
}

func NewBlockchain(conf *utils.Config, to string) *Blockchain {
	bolt_db := db.GetDB(conf)
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

			err = b.Put([]byte(lastBlock), genesis.Hash)
			if err != nil {
				panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte(lastBlock))
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	blockchainIterator := func() BlockchainIterator {
		return &BlockchainIteratorImpl{curHash: tip, db: bolt_db}
	}

	return &Blockchain{TipHash: tip, DB: bolt_db, Iterator: blockchainIterator}
}
