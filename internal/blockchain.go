package internal

import (
	"log"

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
	defer db.Close()

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
	// checking if we need genesis block
	var prevHash []byte
	if len(bc.Blocks) > 0 {
		prevBlock := bc.Blocks[len(bc.Blocks)-1]
		prevHash = prevBlock.Hash
	}

	block := NewBlock(prevHash, tx)

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
