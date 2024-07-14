package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"math"
	"math/big"
	"strconv"
)

const difficulty = 24

func Calculate(block *Block) ([32]byte, int) {
	target := getTarget()

	var hashInInt big.Int
	var hash [32]byte
	nonce := 0

	// Implement maximum for nonce?
	for nonce < math.MaxInt64 {
		// prepare data
		data := prepareData(block, nonce)
		hash = sha256.Sum256(data)
		hashInInt.SetBytes(hash[:])

		if hashInInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return hash, nonce
}

func getTarget() *big.Int {
	target := big.NewInt(1)
	target = target.Lsh(target, uint(256-difficulty))
	return target
}

func prepareData(block *Block, nonce int) []byte {
	data := bytes.Join([][]byte{
		block.PrevHash,
		hashTransactions(block.Transactions),
		[]byte(strconv.FormatInt(block.Timestamp, 10)),
		[]byte(strconv.Itoa(difficulty)),
		[]byte(strconv.Itoa(nonce)),
	}, []byte{})

	return data
}

// TODO: implement proper hashing.
// for now, we just serialize and hash.
func hashTransactions(txs []string) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(txs)
	return buf.Bytes()
}
