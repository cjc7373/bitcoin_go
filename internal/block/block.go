package block

import (
	"context"
	"iter"
	"time"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/common"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
)

func newBlock(ctx context.Context, txs []*block_proto.Transaction, prevBlockHash []byte) *block_proto.Block {
	block := &block_proto.Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  txs,
		PrevBlockHash: prevBlockHash,
		Hash:          nil,
		Nonce:         0,
	}
	pow := NewProofOfWork(block)
	nonce, hash, err := pow.Run(ctx)
	if err != nil {
		return nil
	}

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func GetBlock(db *bolt.DB, hash []byte) (*block_proto.Block, error) {
	block := &block_proto.Block{}
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(common.BlockBucket))
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

// iterate blocks from end to start
func AllBlocks(db *bolt.DB, tipHash []byte) iter.Seq[*block_proto.Block] {
	curHash := tipHash
	return func(yield func(*block_proto.Block) bool) {
		for len(curHash) != 0 {
			block := &block_proto.Block{}
			err := db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(common.BlockBucket))

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

func saveBlock(tx *bolt.Tx, newBlock *block_proto.Block) error {
	b := tx.Bucket([]byte(common.BlockBucket))
	data, err := proto.Marshal(newBlock)
	if err != nil {
		return err
	}

	if err := b.Put(newBlock.Hash, data); err != nil {
		return err
	}
	return nil
}

func AddBlock(ctx context.Context, db *bolt.DB, bc *block_proto.Blockchain, txs []*block_proto.Transaction) (*block_proto.Block, error) {
	block := newBlock(ctx, txs, bc.TipHash)
	err := db.Update(func(tx *bolt.Tx) error {
		if err := saveBlock(tx, block); err != nil {
			return err
		}

		bc.TipHash = block.Hash
		bc.Height += 1
		if err := saveBlockchain(tx, bc); err != nil {
			return err
		}

		return nil
	})

	return block, err
}

func AddGenesisBlock(db *bolt.DB, to wallet.Address) error {
	// genesis block's prevHash is nil
	bc := &block_proto.Blockchain{}
	txs := []*block_proto.Transaction{NewCoinbaseTransaction(to, []byte("Genesis Block"))}
	_, err := AddBlock(context.Background(), db, bc, txs)
	return err
}
