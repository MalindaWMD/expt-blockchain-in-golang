package main

import (
	"github.com/MalindaWMD/expt-blockchain-in-golang/internal"
)

func main() {
	bc := internal.NewBlockchain()
	defer bc.DB.Close()

	// bc.AddBlock([]string{"Sending more money"})

	bc.Print()
}
