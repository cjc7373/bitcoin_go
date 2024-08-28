package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func EncodePubKey(private *ecdsa.PrivateKey) []byte {
	return append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
}

func ParsePubKey(pubkey []byte) *ecdsa.PublicKey {
	curve := elliptic.P256()
	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubkey)
	x.SetBytes(pubkey[:(keyLen / 2)])
	y.SetBytes(pubkey[(keyLen / 2):])

	return &ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
}

type PubKeyHash []byte
type PubKeyHashSized [ripemd160.Size]byte

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

type Empty struct{}

type Set[T comparable] map[T]Empty

func NewSet[T comparable](items ...T) Set[T] {
	ss := make(Set[T], len(items))
	ss.Insert(items...)
	return ss
}

func (s Set[T]) Insert(items ...T) Set[T] {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

func (s Set[T]) Has(item T) bool {
	_, contained := s[item]
	return contained
}
