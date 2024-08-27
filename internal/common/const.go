package common

// in block bucket, we'll have:
// 32-byte block hash -> block data, encoded by protobuf
// "blockchain" -> blockchain metadata
var BlockBucket = []byte("block")

// in transaction bucket, we'll have:
// 32-byte tx hash -> transaction data
var TransactionBucket = []byte("transaction")

// in uxto bucket, we'll have:
// 20-byte pubkey hash -> txhash and the unspent output indexes in that tx
var UTXOBucket = []byte("chainstate")
