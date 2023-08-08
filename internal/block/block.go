package block

import (
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/cjc7373/bitcoin_go/internal/db"
	"github.com/cjc7373/bitcoin_go/internal/utils"
)

const blockBucket = "block"
const lastBlock = "last_block"

type Block struct {
	Timestamp     int64 // a Unix timestamp
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

type Blockchain struct {
	TipHash []byte // top block hash
	Height  int64
	DB      *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	newBlock := NewBlock(data, bc.TipHash)
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

		err = b.Put([]byte("l"), newBlock.Hash)
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

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain(conf *utils.Config) *Blockchain {
	bolt_db := db.GetDB(conf)
	var tip []byte
	var blocks []*Block

	err := bolt_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))

		if b == nil {
			genesis := NewGenesisBlock()
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
			blocks = append(blocks, genesis)
		} else {
			tip = b.Get([]byte(lastBlock))

			var block Block
			// form a fake block
			block.PrevBlockHash = tip
			// TODO: why block.PrevBlockHash != nil doesn't work?
			for len(block.PrevBlockHash) != 0 {
				data := b.Get(block.PrevBlockHash)
				err := json.Unmarshal(data, &block)
				if err != nil {
					panic(err)
				}
				blocks = append(blocks, &block)
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return &Blockchain{TipHash: tip, DB: bolt_db}
}
