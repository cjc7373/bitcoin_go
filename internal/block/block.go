package block

import (
	"encoding/json"
	"iter"
	"log/slog"
	"time"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/utils"
)

const blockBucket = "block"
const lastBlockKey = "last_block"

var logger = slog.Default()

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

func NewBlock(txs []*block_proto.Transaction, prevBlockHash []byte) *block_proto.Block {
	block := &block_proto.Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  txs,
		PrevBlockHash: prevBlockHash,
		Hash:          nil,
		Nonce:         0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

type Blockchain struct {
	TipHash          []byte // top block hash
	Height           int64
	DB               *bolt.DB
	NewBlockIterator func() utils.Iterator[*Block]
}

func GetBlock(db *bolt.DB, hash []byte) (*block_proto.Block, error) {
	block := &block_proto.Block{}
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		v := b.Get(hash)
		if err := proto.Unmarshal(v, block); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return block, nil
}

func AllBlocks(db *bolt.DB, tipHash []byte) iter.Seq[*block_proto.Block] {
	curHash := tipHash
	return func(yield func(*block_proto.Block) bool) {
		for len(curHash) != 0 {
			block := &block_proto.Block{}
			err := db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(blockBucket))

				v := b.Get(curHash)
				if err := proto.Unmarshal(v, block); err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				logger.Error("err iterate chain", "error", err)
				return
			}

			curHash = block.PrevBlockHash
			if !yield(block) {
				return
			}
		}
	}
}

func AddBlock(db *bolt.DB, bc *block_proto.Blockchain, txs []*block_proto.Transaction) {
	newBlock := NewBlock(txs, bc.TipHash)
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		data, err := json.Marshal(newBlock)
		if err != nil {
			panic(err)
		}

		err = b.Put(newBlock.Hash, data)
		if err != nil {
			panic(err)
		}

		err = b.Put([]byte(lastBlockKey), newBlock.Hash)
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

func NewGenesisBlock(to string) *block_proto.Block {
	return NewBlock([]*block_proto.Transaction{NewCoinbaseTransaction(to, []byte("Genesis Block"))}, nil)
}
