package main

import (
	"fmt"
	"log"

	"github.com/MalindaWMD/expt-blockchain-in-golang/internal"
)

func main() {
	bc := internal.NewBlockchain()
	defer bc.DB.Close()

	from1 := "1F24YsuVtXTz6e5d3FvAjxef2d1AbPqVYAoKDhxEiZGdiLxXM5BUKPsN4ZhzcvAvQi7e8XMq"
	from2 := "1R5Xq1eGS46yPb7c5GDgGsABAWMkZFs9QDYKM2c84GS2DcLAfEbnbJimtRkbJ87WTYobFMGm"

	addr1 := internal.GetAddress(from1)
	addr2 := internal.GetAddress(from2)

	addr1Pubkeyhash := internal.GetPubKeyHash(addr1.PublicKey)
	addr2Pubkeyhash := internal.GetPubKeyHash(addr2.PublicKey)

	fmt.Printf("%x ==> %x\n\n", addr1Pubkeyhash, addr2Pubkeyhash)

	// // 1st time only
	// ctx := bc.NewCoinbaseTransaction(addr1Pubkeyhash, addr1.PublicKey)
	// bc.AddBlock([]*internal.Transaction{ctx})

	_, balance := bc.GetUTXOs(addr1Pubkeyhash)
	fmt.Printf("%x balance: %d\n\n", addr1Pubkeyhash, balance)

	// FROM 1 to 2
	tx, err := bc.NewTransaction(addr1Pubkeyhash, addr2Pubkeyhash, 2, addr1.PublicKey)
	if err != nil {
		log.Println("TX:", err)
	}
	if tx != nil {
		fmt.Printf("Adding new transaction: %x\n", tx.ID)
		bc.AddBlock([]*internal.Transaction{tx})
	}

	// FROM 2 to 1
	tx, err = bc.NewTransaction(addr2Pubkeyhash, addr1Pubkeyhash, 3, addr2.PublicKey)
	if err != nil {
		log.Println("TX:", err)
	}
	if tx != nil {
		bc.AddBlock([]*internal.Transaction{tx})
	}

	// TODO:
	// Clean up transactions.go
	// separate addresses
	//

	_, balance = bc.GetUTXOs(addr1Pubkeyhash)
	fmt.Printf("%x balance: %d\n\n", addr1Pubkeyhash, balance)

	_, balance = bc.GetUTXOs(addr2Pubkeyhash)
	fmt.Printf("%x balance after: %d\n\n", addr2Pubkeyhash, balance)

	bc.Print()
}
