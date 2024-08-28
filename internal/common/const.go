package common

// in block bucket, we'll have:
// 32-byte block hash -> block data, encoded by protobuf
// "blockchain" -> blockchain metadata
var BlockBucket = []byte("block")

// transaction bucket, as well as uxto bucket,
// serves as a cache of the blockchain,
// and can be rebuilt if needed

// in transaction bucket, we'll have:
// 32-byte tx hash -> transaction data
var TransactionBucket = []byte("transaction")

type TxHashSized [32]byte

// in uxto bucket, we'll have:
// 20-byte pubkey hash -> utxoSet (txhash and the unspent output indexes in that tx)
var UTXOBucket = []byte("chainstate")
