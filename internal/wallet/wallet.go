package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cjc7373/bitcoin_go/internal/utils"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

const DefatltWalletName = "default"

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

func NewWalletFromPEM(pemEncoded []byte) *Wallet {
	block, _ := pem.Decode(pemEncoded)

	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		panic(err)
	}
	wallet := Wallet{*privateKey, utils.EncodePubKey(privateKey)}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}
	pubKey := utils.EncodePubKey(private)

	return *private, pubKey
}

func (w *Wallet) EncodeToPEM() []byte {
	x509Encoded, err := x509.MarshalECPrivateKey(&w.PrivateKey)
	if err != nil {
		panic(err)
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return pemEncoded

}

func (w *Wallet) GetAddress() string {
	// we try to generate a legal bitcoin address here, so we follow the protocol
	// https://en.bitcoin.it/wiki/Protocol_documentation#Addresses
	pubKeyHash := utils.HashPubKey(w.PublicKey)

	version := []byte{0} // P2PKH address
	versionedPayload := append(version, pubKeyHash...)

	firstSHA := sha256.Sum256(versionedPayload)
	secondSHA := sha256.Sum256(firstSHA[:])
	checksum := secondSHA[:4]

	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)

	return address
}
