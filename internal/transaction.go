package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
)

const reward = 10

type Transaction struct {
	ID      []byte
	Inputs  []*TxInput
	Outputs []*TxOutput
}

type TxInput struct {
	TxId      []byte
	OutIdx    int
	Signature []byte
	PublicKey []byte //optional
}

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// New transaction: pass from, to
//	- get all UTXOs belongs to the sender
// 	- validate balance
//	- if valid, create inputs from outputs from UTXOs
//	- create outputs
//	- return transaction
// TODO: implement mem pool
//	- add transaction to a block
// 	- mine

// NOTE: We'll ust pubkeyhas for now, will implement addresses later
func (bc *Blockchain) NewTransaction(from, to []byte, amount int, fromPubKey []byte) (*Transaction, error) {
	utxos, balance := bc.GetUTXOs(from)

	fmt.Printf("Trx: From: %x => To: %x\n", from, to)

	if amount > balance {
		return nil, errors.New("insufficiant balance")
	}

	var inputs []*TxInput
	var outputs []*TxOutput

	calculatedAmount := 0
	for _, tx := range utxos {
		for idx, o := range tx.Outputs {
			if !bytes.Equal(o.PubKeyHash, from) {
				continue
			}

			inputs = append(inputs, &TxInput{
				TxId:      tx.ID,
				OutIdx:    idx,
				Signature: []byte("ABC123"),
				PublicKey: fromPubKey,
			})

			calculatedAmount += o.Value

			if calculatedAmount >= amount {
				break
			}
		}
	}

	outputs = append(outputs, &TxOutput{
		Value:      amount,
		PubKeyHash: to,
	})

	// "change" transaction
	if calculatedAmount > amount {
		outputs = append(outputs, &TxOutput{
			Value:      calculatedAmount - amount,
			PubKeyHash: from,
		})
	}

	tx := Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}
	tx.SetId()

	return &tx, nil
}

// TODO: TO be used in the block assembling stage. but for now we'll just use it directly.
func (bc *Blockchain) NewCoinbaseTransaction(pubkeyhash []byte, fromPubKey []byte) *Transaction {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("1111111111111111111111111111111111111111111111111111111111111111")
	encoded := buf.Bytes()
	hashed := sha256.Sum256(encoded)
	txId := hashed[:]

	inputs := []*TxInput{
		{
			TxId:      txId,
			OutIdx:    -1,
			PublicKey: fromPubKey,
		},
	}

	outputs := []*TxOutput{
		{
			Value:      reward,
			PubKeyHash: pubkeyhash,
		},
	}
	return &Transaction{
		ID:      txId,
		Inputs:  inputs,
		Outputs: outputs,
	}
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

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}

func (tx *Transaction) SetId() {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(tx)
	encoded := buf.Bytes()
	id := sha256.Sum256(encoded)

	tx.ID = id[:]
}

// func getSpentOutputs(blocks []*Block, fromHash []byte) map[string][]int {
// 	spentOutputs := make(map[string][]int)

// 	for _, block := range blocks {
// 		for _, tx := range block.Transactions {
// 			for _, i := range tx.Inputs {
// 				pubHash := GetPubKeyHash(i.PublicKey)
// 				if !bytes.Equal(pubHash, fromHash) {
// 					continue
// 				}

// 				txId := hex.EncodeToString(i.TxId)
// 				spentOutputs[txId] = append(spentOutputs[txId], i.OutIdx)
// 			}
// 		}
// 	}

// 	return spentOutputs
// }
