package internal

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Blockchain struct {
	Blocks []*Block
}

type Block struct {
	PrevHash     []byte
	Hash         []byte
	Timestamp    int64
	Nonce        int
	Transactions []string
}

func NewBlockchain() *Blockchain {
	// create new blockchain
	bc := &Blockchain{
		Blocks: []*Block{},
	}

	if len(bc.Blocks) == 0 {
		bc.AddBlock([]string{"Genesis block"})
	}

	return bc
}

func (bc *Blockchain) AddBlock(tx []string) {
	// checking if we need genesis block
	var prevHash []byte
	if len(bc.Blocks) > 0 {
		prevBlock := bc.Blocks[len(bc.Blocks)-1]
		prevHash = prevBlock.Hash
	}

	b := NewBlock(prevHash, tx)

	bc.Blocks = append(bc.Blocks, b)
}

func NewBlock(prevHash []byte, tx []string) *Block {
	// Prepare data
	timestamp := time.Now().Unix()
	data := bytes.Join(
		[][]byte{prevHash, []byte(strconv.FormatInt(timestamp, 10))},
		[]byte{},
	)
	hash := sha256.Sum256(data)

	// new block
	b := &Block{
		PrevHash:     prevHash,
		Hash:         hash[:],
		Timestamp:    timestamp,
		Nonce:        0,
		Transactions: tx,
	}

	return b
}
