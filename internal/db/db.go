package db

import (
	"github.com/cjc7373/bitcoin_go/internal/utils"
	bolt "go.etcd.io/bbolt"
)

const BucketName = "blockchain"

func GetDB(conf *utils.Config) *bolt.DB {
	db, err := bolt.Open(conf.DBPath, 0600, nil)
	if err != nil {
		panic(err)
	}

	return db
}

// 'b' + 32-byte block hash -> block data
// "last_block" -> the hash of the last block in a chain
