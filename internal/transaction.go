package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
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
func (bc *Blockchain) NewTransaction(from, to string, amount int) (*Transaction, error) {

	fmt.Println("Sending...")
	fmt.Println("\tFrom \t: ", from)
	fmt.Println("\tTo \t: ", to)
	fmt.Println("\tAmount \t: ", amount)
	fmt.Println()

	fromAddr := GetAddress(from)
	fromPubKeyHash := GetPubKeyHash(fromAddr.PublicKey)

	utxos, balance := bc.GetUTXOs(fromPubKeyHash)

	if amount > balance {
		return nil, errors.New("insufficiant balance")
	}

	var inputs []*TxInput
	var outputs []*TxOutput

	calculatedAmount := 0
	for _, tx := range utxos {
		for idx, o := range tx.Outputs {
			if !bytes.Equal(o.PubKeyHash, fromPubKeyHash) {
				continue
			}

			inputs = append(inputs, &TxInput{
				TxId:      tx.ID,
				OutIdx:    idx,
				Signature: []byte("ABC123"),
				PublicKey: fromAddr.PublicKey,
			})

			calculatedAmount += o.Value

			if calculatedAmount >= amount {
				break
			}
		}
	}

	toAddr := GetAddress(to)
	outputs = append(outputs, &TxOutput{
		Value:      amount,
		PubKeyHash: GetPubKeyHash(toAddr.PublicKey),
	})

	// "change" transaction
	if calculatedAmount > amount {
		outputs = append(outputs, &TxOutput{
			Value:      calculatedAmount - amount,
			PubKeyHash: fromPubKeyHash,
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
func (bc *Blockchain) NewCoinbaseTransaction(from string) *Transaction {
	fromAddr := GetAddress(from)
	fromPubKeyHash := GetPubKeyHash(fromAddr.PublicKey)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("1")
	encoded := buf.Bytes()
	hashed := sha256.Sum256(encoded)
	txId := hashed[:]

	inputs := []*TxInput{
		{
			TxId:      txId,
			OutIdx:    -1,
			PublicKey: fromAddr.PublicKey,
		},
	}

	outputs := []*TxOutput{
		{
			Value:      reward,
			PubKeyHash: fromPubKeyHash,
		},
	}
	return &Transaction{
		ID:      txId,
		Inputs:  inputs,
		Outputs: outputs,
	}
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
