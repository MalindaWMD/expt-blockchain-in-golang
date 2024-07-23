package internal

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"log"

	"github.com/boltdb/bolt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressBucket = "addresses"

type Address struct {
	Address    string
	PrivateKey []byte
	PublicKey  []byte
}

// TODO: Address generation should be moved to a wallet.
// generate address
func GetAddress(address string) *Address {
	db, err := bolt.Open("./internal/db/addresses.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer db.Close()

	var addr *Address
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(addressBucket))
		if err != nil {
			return err
		}

		entry := b.Get([]byte(address))
		if entry == nil {
			// TODO: Encode private key
			priv, pub := getKeyPair()
			ad := generateAddress(pub)
			addr = &Address{
				Address:    ad,
				PrivateKey: priv,
				PublicKey:  pub,
			}
			b.Put([]byte(address), addr.serialize())
			return nil
		}

		addr = deserialize(entry)

		return nil
	})
	if err != nil {
		log.Fatal("Error fetching address.", err)
	}

	return addr
}

func generateAddress(pub []byte) string {
	pubKeyHash := GetPubKeyHash(pub)
	checksum := getChecksum(pubKeyHash)

	payload := bytes.Join([][]byte{
		{version},
		pubKeyHash,
		checksum[:],
	}, []byte{})

	return base58.Encode(payload)
}

func GetPubKeyHash(pub []byte) []byte {
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

func getKeyPair() ([]byte, []byte) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal("Error generating private key.", err)
	}

	publicKey := privateKey.PublicKey

	priv, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	pub, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		log.Fatal(err)
	}

	return priv, pub
}

func (a *Address) serialize() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(a)

	return buf.Bytes()
}

func deserialize(data []byte) *Address {
	var buf bytes.Buffer
	var addr *Address
	buf.Write(data)

	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&addr)
	if err != nil {
		return nil
	}

	return addr
}
