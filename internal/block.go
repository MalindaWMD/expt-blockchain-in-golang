package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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

	return b
}

func GenesisBlock() *Block {
	tx := []*Transaction{}
	genesis := NewBlock([]byte{}, tx)
	hashData := genesis.PrepareData(genesis.Nonce)

	hash := sha256.Sum256(hashData)

	genesis.Hash = hash[:]

	return genesis
}

func (b *Block) Mine() *Block {
	hash, nonce := Calculate(b)
	b.Hash = hash[:]
	b.Nonce = nonce

	return b
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
