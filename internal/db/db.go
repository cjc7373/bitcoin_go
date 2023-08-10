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
