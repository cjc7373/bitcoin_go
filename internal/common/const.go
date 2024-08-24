package common

// in block bucket, we'll have:
// 32-byte block hash -> block data, encoded by json
// "last_block" -> the hash of the last block in a chain
var BlockBucket = []byte("block")

var UTXOBucket = []byte("chainstate")
