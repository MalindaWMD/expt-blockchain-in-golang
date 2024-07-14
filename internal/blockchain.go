package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

const dbPath = "./internal/db/blockchain.db"
const bucketName = "blocks"

type Blockchain struct {
	DB     *bolt.DB
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var blocks []*Block

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		c := b.Cursor()

		for k, blockData := c.First(); k != nil; k, blockData = c.Next() {
			block := DeeserializeBlockData(blockData)
			blocks = append(blocks, block)
		}

		return nil
	})
	if err != nil {
		log.Fatal("Error reading DB:", err)
	}

	// create new blockchain
	bc := &Blockchain{
		DB:     db,
		Blocks: blocks,
	}

	// No blocks, add genesis block.
	if len(blocks) == 0 {
		log.Println("Adding genesis block")
		bc.AddBlock([]string{"Genesis block tx"})
	}

	return bc
}

func (bc *Blockchain) AddBlock(tx []string) {
	log.Println("Adding new block.")
	// checking if we need genesis block
	var prevHash []byte
	if len(bc.Blocks) > 0 {
		prevBlock := bc.Blocks[len(bc.Blocks)-1]
		prevHash = prevBlock.Hash
	}

	block := NewBlock(prevHash, tx)

	log.Println("Calculating PoW")
	// run PoW
	hash, nonce := Calculate(block)
	block.Hash = hash[:]
	block.Nonce = nonce

	log.Println("Updating db")
	bc.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		blockData := SerializeBlock(block)
		b.Put(block.Hash, blockData)

		return nil
	})

	bc.Blocks = append(bc.Blocks, block)
}

func (bc *Blockchain) Print() {
	for _, b := range bc.Blocks {
		tm := time.Unix(b.Timestamp, 0)
		fmt.Printf("Time\t\t: %s\n", tm)
		fmt.Printf("Prev. Hash\t: %x\n", b.PrevHash)
		fmt.Printf("Hash\t\t: %x\n", b.Hash)
		fmt.Printf("Txs\t\t: %s\n", b.Transactions)
		fmt.Printf("Nonce\t\t: %d\n\n", b.Nonce)
	}
}
