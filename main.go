package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/MalindaWMD/expt-blockchain-in-golang/internal"
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

func main() {
	bc := internal.NewBlockchain()
	defer bc.DB.Close()

	from1 := "1F24YsuVtXTz6e5d3FvAjxef2d1AbPqVYAoKDhxEiZGdiLxXM5BUKPsN4ZhzcvAvQi7e8XMq"
	from2 := "1R5Xq1eGS46yPb7c5GDgGsABAWMkZFs9QDYKM2c84GS2DcLAfEbnbJimtRkbJ87WTYobFMGm"

	addr1 := getAddress(from1)
	addr2 := getAddress(from2)

	addr1Pubkeyhash := getPubKeyHash(addr1.PublicKey)
	addr2Pubkeyhash := getPubKeyHash(addr2.PublicKey)

	fmt.Printf("%x ==> %x\n\n", addr1Pubkeyhash, addr2Pubkeyhash)

	// // 1st time only
	// ctx := bc.NewCoinbaseTransaction(addr1Pubkeyhash, addr1.PublicKey)
	// bc.AddBlock([]*internal.Transaction{ctx})

	_, balance := bc.GetUTXOs(addr1Pubkeyhash)
	fmt.Printf("%x balance: %d\n\n", addr1Pubkeyhash, balance)

	tx, err := bc.NewTransaction(addr1Pubkeyhash, addr2Pubkeyhash, 2, addr1.PublicKey)
	if err != nil {
		log.Println("TX:", err)
	}
	if tx != nil {
		fmt.Printf("Adding new transaction: %x\n", tx.ID)
		bc.AddBlock([]*internal.Transaction{tx})
	}

	// TODO:
	// Let's keep a separate list of UTXOs in a bucket and update that accordingly?
	//

	_, balance = bc.GetUTXOs(addr1Pubkeyhash)
	fmt.Printf("%x balance: %d\n\n", addr1Pubkeyhash, balance)

	_, balance = bc.GetUTXOs(addr2Pubkeyhash)
	fmt.Printf("%x balance after: %d\n\n", addr2Pubkeyhash, balance)

	/* Addresses
	1F24YsuVtXTz6e5d3FvAjxef2d1AbPqVYAoKDhxEiZGdiLxXM5BUKPsN4ZhzcvAvQi7e8XMq
	1R5Xq1eGS46yPb7c5GDgGsABAWMkZFs9QDYKM2c84GS2DcLAfEbnbJimtRkbJ87WTYobFMGm
	*/

	// addr1 := getAddress(bc.DB, "1F24YsuVtXTz6e5d3FvAjxef2d1AbPqVYAoKDhxEiZGdiLxXM5BUKPsN4ZhzcvAvQi7e8XMq")
	// log.Println("Addr:", addr1.Address)
	// bc.AddBlock([]string{"Sending 1 from: " + addr1.Address})

	// addr2 := getAddress(bc.DB, "")
	// log.Println("Addr:", addr2.Address)
	// bc.AddBlock([]string{"Sending 3 from: " + addr2.Address})

	bc.Print()
}

// TODO: Address generation should be moved to a wallet.
// generate address
func getAddress(address string) *Address {
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
	pubKeyHash := getPubKeyHash(pub)
	checksum := getChecksum(pubKeyHash)

	payload := bytes.Join([][]byte{
		{version},
		pubKeyHash,
		checksum[:],
	}, []byte{})

	return base58.Encode(payload)
}

func getPubKeyHash(pub []byte) []byte {
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
