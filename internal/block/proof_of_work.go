package block

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/utils"
)

const maxNonce = math.MaxInt64

// The difficulty to mine
// It's the leading zeros of the result hash
// like 0x000000abcd...
// We won’t implement a target adjusting algorithm for simplicity
var targetBits = 16

// in nodes we will deliberately add some latency to make pow slow
// so that it won't actually consume much CPU time
var PowSleepTime time.Duration = 0

type ProofOfWork struct {
	block  *block_proto.Block
	target *big.Int
}

func NewProofOfWork(b *block_proto.Block) *ProofOfWork {
	target := big.NewInt(1)
	// hex of target will be 0x00000100000....
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	// TODO
	txData, err := json.Marshal(pow.block.Transactions)
	if err != nil {
		panic(err)
	}
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			txData,
			utils.IntToHex(pow.block.Timestamp),
			utils.IntToHex(int64(targetBits)),
			utils.IntToHex(nonce),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) setNonce(data []byte, nonce int64) []byte {
	trimmedData := data[:len(data)-8]
	return append(trimmedData, utils.IntToHex(nonce)...)
}

func (pow *ProofOfWork) Run(ctx context.Context) (int64, []byte, error) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64

	logger.Info("Mining started")
	for _, tx := range pow.block.Transactions {
		logger.Info("With transaction", "tx", (*block_proto.TransactionPretty)(tx))
	}
	data := pow.prepareData(nonce)
	batchSize := 10000
outer:
	for nonce < maxNonce {
		cnt := 0
		// FIXME: select turns out to be very slow, find out why
		select {
		case <-ctx.Done():
			logger.Info("Mining stopped", "err", ctx.Err())
			return 0, nil, ctx.Err()
		default:
			for cnt < batchSize {
				cnt++
				data := pow.setNonce(data, nonce)
				// we won't use key derivation functions like PBKDF2 and scrypt for simplicity
				hash = sha256.Sum256(data)
				hashInt.SetBytes(hash[:])

				if hashInt.Cmp(pow.target) == -1 {
					break outer
				} else {
					nonce++
				}

			}
			time.Sleep(PowSleepTime)
		}
	}
	logger.Info("mining completed", "hash", fmt.Sprintf("%x", hash), "nonce", nonce)

	return nonce, hash[:], nil
}
