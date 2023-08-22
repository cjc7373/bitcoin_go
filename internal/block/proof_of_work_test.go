package block

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newDumbBlock() *Block {
	return &Block{time.Now().Unix(), []Transaction{}, []byte{}, []byte{}, 0}
}

func BenchmarkPrepareData(b *testing.B) {
	block := newDumbBlock()
	pow := NewProofOfWork(block)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pow.prepareData(1)
	}
}

func BenchmarkSetNonce(b *testing.B) {
	block := newDumbBlock()
	pow := NewProofOfWork(block)
	data := pow.prepareData(0)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		data = pow.setNonce(data, 1)
	}
}

func TestPoW(t *testing.T) {
	block := newDumbBlock()
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	data := pow.prepareData(nonce)
	actualHash := sha256.Sum256(data)
	assert.Equal(t, hash, actualHash[:])
}
