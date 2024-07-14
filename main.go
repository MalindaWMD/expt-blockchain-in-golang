package main

import (
	"fmt"
	"time"

	"github.com/MalindaWMD/expt-blockchain-in-golang/internal"
)

func main() {
	bc := internal.NewBlockchain()

	for _, b := range bc.Blocks {
		tm := time.Unix(b.Timestamp, 0)
		fmt.Printf("Time\t\t: %s\n", tm)
		fmt.Printf("Prev. Hash\t: %x\n", b.PrevHash)
		fmt.Printf("Hash\t\t: %x\n", b.Hash)
		fmt.Printf("Txs\t\t: %s\n\n", b.Transactions)
	}
}
