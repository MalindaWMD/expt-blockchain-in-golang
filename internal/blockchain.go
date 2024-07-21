package internal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

const dbPath = "./internal/db/blockchain.db"
const bucketName = "blocks"
const metadataBucket = "metadata"

type Blockchain struct {
	DB   *bolt.DB
	Tip  []byte
	lock sync.Mutex
}

type Itarator struct {
	DB          *bolt.DB
	CurrentHash []byte
}

func NewBlockchain() *Blockchain {
	// check if DB file exists
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		return initBlockchain()
	}

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(metadataBucket))
		tip = b.Get([]byte("latest"))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	bc := &Blockchain{
		DB:  db,
		Tip: tip,
	}

	return bc
}

func initBlockchain() *Blockchain {
	// create db and buckets
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	log.Println("Creating genesis block")
	genesis := GenesisBlock()

	err = db.Update(func(tx *bolt.Tx) error {
		metaBucket, err := tx.CreateBucket([]byte(metadataBucket))
		if err != nil {
			return err
		}

		blocksBucket, err := tx.CreateBucket([]byte(bucketName))
		if err != nil {
			return err
		}

		blocksBucket.Put(genesis.Hash, genesis.SerializeBlock())
		metaBucket.Put([]byte("latest"), genesis.Hash)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	bc := &Blockchain{
		DB:  db,
		Tip: genesis.Hash,
	}

	return bc
}

func (bc *Blockchain) AddBlock(tx []*Transaction) *Block {
	log.Println("Adding new block.")

	block := NewBlock(bc.Tip, tx)

	// Just calling Validate() here for now.
	// TODO: implement it when a block PoW should be validated after broadcasting it.
	log.Println("Block validity:", Validate(block))

	// TODO: Mining process should be separated
	log.Println("Mining...")
	block = block.Mine()

	log.Println("Updating db")
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		blockData := block.SerializeBlock()
		b.Put(block.Hash, blockData)

		// We need to keep the latest hash we have added in order to set the tip.
		// Cannot get the last one since BoltDB is byte-sorted.
		b = tx.Bucket([]byte(metadataBucket))
		b.Put([]byte("latest"), block.Hash)

		return nil
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	bc.Tip = block.Hash

	return block
}

func (bc *Blockchain) NewItarator() *Itarator {
	return &Itarator{
		DB:          bc.DB,
		CurrentHash: bc.Tip,
	}
}

// TODO: Implement a cursor that'll get blocks in added order.
func (bc *Blockchain) Blocks() []*Block {
	var blocks []*Block

	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			block := DeserializeBlockData(v)
			blocks = append(blocks, block)
		}

		return nil
	})

	return blocks
}

func (i *Itarator) Next() *Block {
	var block *Block
	i.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		value := b.Get([]byte(i.CurrentHash))

		block = DeserializeBlockData(value)

		return nil
	})

	if block != nil {
		i.CurrentHash = block.PrevHash
	}

	return block
}

func (bc *Blockchain) Print() {
	i := bc.NewItarator()

	for {
		b := i.Next()
		if b == nil {
			break
		}

		tm := time.Unix(b.Timestamp, 0)
		fmt.Printf("Time\t\t: %s\n", tm)
		fmt.Printf("Prev. Hash\t: %x\n", b.PrevHash)
		fmt.Printf("Hash\t\t: %x\n", b.Hash)
		fmt.Printf("Txs\t\t: %v\n", b.Transactions)
		fmt.Printf("Nonce\t\t: %d\n\n", b.Nonce)
	}
}
