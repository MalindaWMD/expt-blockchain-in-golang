package internal

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
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

	utxos, balance := bc.GetUTXOs(fromAddr)

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
	toPubKeyhash := GetPubKeyHash(toAddr.PublicKey)

	fmt.Printf("To Public KeyHash: %x \n\n", toPubKeyhash)

	outputs = append(outputs, &TxOutput{
		Value:      amount,
		PubKeyHash: toPubKeyhash,
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
	tx.Sign(&fromAddr.PrivateKey)

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
	hash := tx.Hash()
	tx.ID = hash[:]
}

func (tx *Transaction) StringId() string {
	return hex.EncodeToString(tx.ID)
}

func (tx *Transaction) Serialize() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(tx)
	return buf.Bytes()
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) *Transaction {
	trimmed := tx.Trim()
	hash := trimmed.Hash()

	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hash)
	if err != nil {
		log.Fatal(err)
	}

	tx.Inputs[0].Signature = signature

	return tx
}

func (tx *Transaction) Verify(publicKey ecdsa.PublicKey) bool {
	if tx.IsCoinbase() {
		return true
	}

	trimmed := tx.Trim()
	hash := trimmed.Hash()

	signature := tx.Inputs[0].Signature

	return ecdsa.VerifyASN1(&publicKey, hash, signature)
}

func (tx *Transaction) Hash() []byte {
	serialized := tx.Serialize()
	hash := sha256.Sum256(serialized)
	return hash[:]
}

func (tx *Transaction) Trim() *Transaction {
	var inputs []*TxInput
	var outputs []*TxOutput

	for _, i := range tx.Inputs {
		ii := TxInput{
			TxId:   i.TxId,
			OutIdx: i.OutIdx,
		}
		inputs = append(inputs, &ii)
	}

	for _, o := range tx.Outputs {
		outputs = append(outputs, &TxOutput{o.Value, o.PubKeyHash})
	}

	return &Transaction{
		ID:      tx.ID,
		Inputs:  inputs,
		Outputs: outputs,
	}
}
