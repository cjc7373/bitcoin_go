package db

import (
	"path"

	"github.com/cjc7373/bitcoin_go/internal/common"
	"github.com/cjc7373/bitcoin_go/internal/utils"
	bolt "go.etcd.io/bbolt"
)

const BucketName = "blockchain"

func OpenDB(conf *utils.Config) *bolt.DB {
	db, err := bolt.Open(path.Join(conf.GetDataDir(), conf.DBPath), 0600, nil)
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(common.BlockBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(common.UTXOBucket); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return db
}
