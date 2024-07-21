package internal

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
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

		hash := bc.Tip

		for {
			blockData := b.Get(hash)
			if blockData == nil {
				break
			}

			block := DeserializeBlockData(blockData)
			blocks = append(blocks, block)

			hash = block.PrevHash
		}

		return nil
	})

	return blocks
}

func (bc *Blockchain) GetUTXOs(pubkeyhash []byte) ([]*Transaction, int) {
	spentOutputs := make(map[string][]int)
	unspentTransactions := []*Transaction{}
	balance := 0

	blocks := bc.Blocks()
	// spentOutputs := getSpentOutputs(blocks, pubkeyhash)

	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID)

			for _, i := range tx.Inputs {
				pubHash := GetPubKeyHash(i.PublicKey)
				if !bytes.Equal(pubHash, pubkeyhash) {
					continue
				}

				txId := hex.EncodeToString(i.TxId)
				spentOutputs[txId] = append(spentOutputs[txId], i.OutIdx)
			}

			for idx, o := range tx.Outputs {
				if !slices.Contains(spentOutputs[txId], idx) {
					// check if the output belongs to the sender.
					// TODO: May be better if we filter this in getSpentOutputs()???
					if !bytes.Equal(o.PubKeyHash, pubkeyhash) {
						continue
					}
					unspentTransactions = append(unspentTransactions, tx)
					balance += o.Value
				}
			}
		}
	}

	return unspentTransactions, balance
}

func (bc *Blockchain) GetBalance(from string) int {
	address := GetAddress(from)
	hash := GetPubKeyHash(address.PublicKey)
	_, balance := bc.GetUTXOs(hash)

	return balance
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
	blocks := bc.Blocks()

	for _, b := range blocks {
		tm := time.Unix(b.Timestamp, 0)
		fmt.Printf("Time\t\t: %s\n", tm)
		fmt.Printf("Prev. Hash\t: %x\n", b.PrevHash)
		fmt.Printf("Hash\t\t: %x\n", b.Hash)
		fmt.Printf("Txs\t\t: %v\n", b.Transactions)
		fmt.Printf("Nonce\t\t: %d\n\n", b.Nonce)
	}
}
