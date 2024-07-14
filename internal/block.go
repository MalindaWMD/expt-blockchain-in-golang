package internal

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	PrevHash     []byte
	Hash         []byte
	Timestamp    int64
	Nonce        int
	Transactions []string
}

func NewBlock(prevHash []byte, tx []string) *Block {
	// Prepare data
	timestamp := time.Now().Unix()

	// new block
	b := &Block{
		PrevHash:     prevHash,
		Timestamp:    timestamp,
		Nonce:        0,
		Transactions: tx,
	}

	return b
}

func SerializeBlock(b *Block) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(b)
	if err != nil {
		return nil
	}

	return buf.Bytes()
}

// TODO: re-use buffer???
func DeeserializeBlockData(data []byte) *Block {
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
