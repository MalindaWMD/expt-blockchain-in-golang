package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/MalindaWMD/expt-blockchain-in-golang/internal"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)

func main() {
	bc := internal.NewBlockchain()
	defer bc.DB.Close()

	// addr1 := getAddress()
	// bc.AddBlock([]string{"Sending 1 from: " + addr1})

	addr2 := getAddress()
	bc.AddBlock([]string{"Sending 3 from: " + addr2})

	bc.Print()
}

// TODO: Address generation should be moved to a wallet.
// generate address
func getAddress() string {
	pubKeyHash := getPubKeyHash()
	checksum := getChecksum(pubKeyHash)

	payload := bytes.Join([][]byte{
		{version},
		pubKeyHash,
		checksum[:],
	}, []byte{})

	return base58.Encode(payload)
}

func getPubKeyHash() []byte {
	_, pub := getKeyPair()

	shaHash := sha256.Sum256(pub)
	ripemd := ripemd160.New()
	_, err := ripemd.Write(shaHash[:])
	if err != nil {
		log.Fatal(err)
	}

	return ripemd.Sum(nil)
}

func getChecksum(pubKeyHash []byte) []byte {
	first := sha256.Sum256(pubKeyHash)
	checksum := sha256.Sum256(first[:])
	return checksum[:]
}

func getKeyPair() (*ecdsa.PrivateKey, []byte) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal("Error generating private key.", err)
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return privateKey, publicKey
}
