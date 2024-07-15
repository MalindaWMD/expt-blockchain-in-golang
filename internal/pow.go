package internal

import (
	"crypto/sha256"
	"math"
	"math/big"
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
		data := block.PrepareData(nonce)
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

func Validate(block *Block) bool {
	var hash [32]byte
	var hashInInt big.Int
	data := block.PrepareData(block.Nonce)
	hash = sha256.Sum256(data)
	hashInInt.SetBytes(hash[:])

	target := getTarget()

	return hashInInt.Cmp(target) == -1
}

func getTarget() *big.Int {
	target := big.NewInt(1)
	target = target.Lsh(target, uint(256-difficulty))
	return target
}
