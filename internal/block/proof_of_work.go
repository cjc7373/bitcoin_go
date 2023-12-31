package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/cjc7373/bitcoin_go/internal/utils"
)

// The difficulty to mine
// It's the leading zeros of the result hash
// like 0x000000abcd...
// We won’t implement a target adjusting algorithm for simplicity
const targetBits = 16
const maxNonce = math.MaxInt64

// in nodes we will deliberately add some latency to make pow slow
// so that it won't actually consume much CPU time
var PowSleepTime time.Duration = 0

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// hex of target will be 0x00000100000....
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
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
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) setNonce(data []byte, nonce int) []byte {
	time.Sleep(PowSleepTime)
	trimmedData := data[:len(data)-8]
	return append(trimmedData, utils.IntToHex(int64(nonce))...)
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing the following txs:\n")
	for _, tx := range pow.block.Transactions {
		fmt.Println(&tx)
	}
	data := pow.prepareData(nonce)
	for nonce < maxNonce {
		data := pow.setNonce(data, nonce)
		// we won't use key derivation functions like PBKDF2 and scrypt for simplicity
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("%x", hash)
	fmt.Print("\n\n")

	return nonce, hash[:]
}
