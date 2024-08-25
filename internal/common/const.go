package common

// in block bucket, we'll have:
// 32-byte block hash -> block data, encoded by json
// "last_block" -> the hash of the last block in a chain
var BlockBucket = []byte("block")

// in uxto bucket, we'll have:
// 32-byte tx hash -> []int, stores unspent output indexes in that tx
var UTXOBucket = []byte("chainstate")
