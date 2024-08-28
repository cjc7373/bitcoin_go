package block

import (
	"context"
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
		nonce, hash, err := pow.Run(context.Background())
		Expect(err).NotTo(HaveOccurred())
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
	// targetBit = 16
	// sleeptime: 0s, elapsed: 24.015026ms
	// sleeptime: 10ms, elapsed: 134.549324ms
	// sleeptime: 100ms, elapsed: 764.407801ms
	//
	// targetBit = 24
	// sleeptime: 0s, elapsed: 3.747273548s
	// sleeptime: 100µs, elapsed: 19.122343777s
	// sleeptime: 1ms, elapsed: 12.675399572s
	// sleeptime: 10ms, elapsed: 44.332263791s
	// sleeptime: 100ms, elapsed: 1m58.866277474s
	XIt("calculates slow down pow time", func() {
		d := []time.Duration{0, time.Microsecond * 100, time.Millisecond, time.Millisecond * 10, time.Millisecond * 100}
		for _, dur := range d {
			PowSleepTime = dur
			targetBits = 24
			var all time.Duration
			cnt := 10
			for i := 0; i < cnt; i++ {
				start := time.Now()
				block := newDumbBlock()
				pow := NewProofOfWork(block)
				_, _, err := pow.Run(context.Background())
				Expect(err).NotTo(HaveOccurred())
				elapsed := time.Since(start)
				all += elapsed
			}
			GinkgoWriter.Printf("=======\nsleeptime: %v, elapsed: %v\n========\n", dur, all/time.Duration(cnt))
		}
	})
})
