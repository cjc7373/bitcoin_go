package block

import (
	"crypto/sha256"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	block_proto "github.com/cjc7373/bitcoin_go/internal/block/proto"
	"github.com/cjc7373/bitcoin_go/internal/wallet"
)

func newDumbBlock() *block_proto.Block {
	w1 := wallet.NewWallet()
	return &block_proto.Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  []*block_proto.Transaction{NewCoinbaseTransaction(w1.GetAddress(), nil)},
		PrevBlockHash: []byte{},
		Hash:          []byte{},
		Nonce:         0,
	}
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

var _ = Describe("pow test", func() {
	It("generates correct hash", func() {
		block := newDumbBlock()
		pow := NewProofOfWork(block)
		nonce, hash := pow.Run()
		data := pow.prepareData(nonce)
		actualHash := sha256.Sum256(data)
		Expect(hash).To(BeEquivalentTo(actualHash))

		By("hash has correct target bits")
		targetBytes := targetBits / 8
		for _, b := range hash[:targetBytes] {
			Expect(b).To(Equal(byte(0)))
		}
	})

	// on my computer, the result is
	// sleeptime: 0s, elapsed: 1.514541ms
	// sleeptime: 1µs, elapsed: 2.816881144s
	// sleeptime: 10µs, elapsed: 1m34.779183894s
	XIt("calculates slow down pow time", func() {
		for _, dur := range []time.Duration{0, time.Microsecond, time.Microsecond * 10} {
			PowSleepTime = dur
			start := time.Now()
			block := newDumbBlock()
			pow := NewProofOfWork(block)
			_, _ = pow.Run()
			elapsed := time.Since(start)
			GinkgoWriter.Printf("sleeptime: %v, elapsed: %v\n", dur, elapsed)
		}
	})
})
