package internal

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/hex"
	"io"
	"log"

	"github.com/boltdb/bolt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressBucket = "addresses"
const encKey = "12345678123456781234567812345678"

type Address struct {
	Address    string
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type DBAddress struct {
	Address    string
	PrivateKey []byte
	PublicKey  []byte
}

type Wallet struct {
	Addresses map[string]Address
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
			b.Put([]byte(ad), addr.serialize())
			return nil
		}

		addr = deserialize(entry)

		return nil
	})
	if err != nil {
		log.Fatal("Error fetching address.", err)
	}

	// fmt.Println("Address: ", addr.Address)

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

func getKeyPair() (ecdsa.PrivateKey, []byte) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal("Error generating private key.", err)
	}

	publicKey := privateKey.PublicKey

	pub, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		log.Fatal(err)
	}

	return *privateKey, pub
}

func (a *Address) serialize() []byte {

	encryptedPrivateKey := EncryptPrivateKey(a.PrivateKey)
	dbAddr := DBAddress{
		Address:    a.Address,
		PrivateKey: []byte(encryptedPrivateKey),
		PublicKey:  a.PublicKey,
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(dbAddr)

	return buf.Bytes()
}

func deserialize(data []byte) *Address {
	var buf bytes.Buffer
	var dbAddr *DBAddress
	buf.Write(data)

	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&dbAddr)
	if err != nil {
		return nil
	}

	decryptedPrivateKey := DecryptPrivateKey(string(dbAddr.PrivateKey))
	addr := &Address{
		Address:    dbAddr.Address,
		PrivateKey: *decryptedPrivateKey,
		PublicKey:  dbAddr.PublicKey,
	}

	return addr
}

func EncryptPrivateKey(key ecdsa.PrivateKey) string {
	// marshan
	priv, err := x509.MarshalECPrivateKey(&key)
	if err != nil {
		log.Fatal(err)
	}

	//encrypt
	cb, err := aes.NewCipher([]byte(encKey))
	if err != nil {
		log.Fatal(err)
	}

	/* NOTES:
	Use GCM mode since AES directly doesn't support arbitrary lengths.
	*/
	gcm, err := cipher.NewGCM(cb)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, priv, nil)

	enc := hex.EncodeToString(ciphertext)

	return enc
}

func DecryptPrivateKey(encrypted string) *ecdsa.PrivateKey {
	cb, err := aes.NewCipher([]byte(encKey))
	if err != nil {
		log.Fatal(err)
	}

	decodedCipherText, err := hex.DecodeString(encrypted)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(cb)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	decryptedData, err := gcm.Open(nil, decodedCipherText[:gcm.NonceSize()], decodedCipherText[gcm.NonceSize():], nil)
	if err != nil {
		log.Fatal(err)
	}

	priv, err := x509.ParseECPrivateKey(decryptedData)
	if err != nil {
		log.Fatal(err)
	}

	return priv
}
