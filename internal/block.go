package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	PrevHash     []byte
	Hash         []byte
	Timestamp    int64
	Nonce        int
	Transactions []*Transaction
}

func NewBlock(prevHash []byte, txs []*Transaction) *Block {
	// Prepare data
	timestamp := time.Now().Unix()

	// new block
	b := &Block{
		PrevHash:     prevHash,
		Timestamp:    timestamp,
		Nonce:        0,
		Transactions: txs,
	}
	b.Hash = b.GetHash()

	return b
}

func GenesisBlock() *Block {
	tx := []*Transaction{}
	genesis := NewBlock([]byte{}, tx)
	data := genesis.PrepareData(genesis.Nonce)
	hash := sha256.Sum256(data)
	genesis.Hash = hash[:]

	return genesis
}

func (b *Block) Mine() *Block {
	if !b.Validate() {
		log.Fatal("Block is invalid.")
		// TODO: Handle this properly
		// return
	}

	hash, nonce := Calculate(b)
	b.Hash = hash[:]
	b.Nonce = nonce

	return b
}

func (b *Block) GetHash() []byte {
	if len(b.Transactions) == 0 {
		return nil
	}

	var rootHash []byte
	hashes := [][]byte{}

	for _, tx := range b.Transactions {
		hashes = append(hashes, tx.ID)
	}

	for {
		if len(hashes)%2 == 1 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}

		subHashes := [][]byte{}
		for i := 0; i < len(hashes); i += 2 {
			hash1 := hashes[i]
			hash2 := hashes[i+1]

			joined := bytes.Join([][]byte{hash1, hash2}, []byte{})
			joinedHash := sha256.Sum256(joined)
			subHashes = append(subHashes, joinedHash[:])
		}

		hashes = subHashes

		if len(hashes) == 1 {
			rootHash = hashes[0]
			break
		}
	}

	return rootHash
}

func (b *Block) Validate() bool {
	newHash := b.GetHash()
	return bytes.Equal(b.Hash, newHash)
}

func (b *Block) PrepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		b.PrevHash,
		b.HashTransactions(),
		[]byte(strconv.FormatInt(b.Timestamp, 10)),
		[]byte(strconv.Itoa(difficulty)),
		[]byte(strconv.Itoa(nonce)),
	}, []byte{})

	return data
}

// TODO: implement proper hashing.
// for now, we just serialize and hash.
func (b *Block) HashTransactions() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(b.Transactions)
	return buf.Bytes()
}

func (b *Block) SerializeBlock() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(b)
	if err != nil {
		return nil
	}

	return buf.Bytes()
}

// TODO: re-use buffer???
func DeserializeBlockData(data []byte) *Block {
	var buf bytes.Buffer
	buf.Write(data)

	var block Block

	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&block)
	if err != nil {
		return nil
	}

	return &block
}
