package main

import (
	"fmt"
	"log"

	"github.com/MalindaWMD/expt-blockchain-in-golang/internal"
)

func main() {
	bc := internal.NewBlockchain()
	defer bc.DB.Close()

	from := "1F24YsuVtXTz6e5d3FvAjxef2d1AbPqVYAoKDhxEiZGdiLxXM5BUKPsN4ZhzcvAvQi7e8XMq"
	to := "1R5Xq1eGS46yPb7c5GDgGsABAWMkZFs9QDYKM2c84GS2DcLAfEbnbJimtRkbJ87WTYobFMGm"

	// // 1st time only
	// ctx := bc.NewCoinbaseTransaction(addr1Pubkeyhash, addr1.PublicKey)
	// bc.AddBlock([]*internal.Transaction{ctx})

	fmt.Printf("\n%s balance: %d\n", from, bc.GetBalance(from))
	fmt.Printf("%s balance: %d\n\n", to, bc.GetBalance(to))

	// FROM 1 to 2
	tx, err := bc.NewTransaction(from, to, 2)
	if err != nil {
		log.Println("TX:", err)
	}
	if tx != nil {
		bc.AddBlock([]*internal.Transaction{tx})
	}

	// // FROM 2 to 1
	// tx, err = bc.NewTransaction(to, from, 3)
	// if err != nil {
	// 	log.Println("TX:", err)
	// }
	// if tx != nil {
	// 	bc.AddBlock([]*internal.Transaction{tx})
	// }

	// TODO:
	// Clean up transactions.go
	// separate addresses
	//

	fmt.Printf("\n%s balance: %d\n", from, bc.GetBalance(from))
	fmt.Printf("%s balance: %d\n\n", to, bc.GetBalance(to))

	bc.Print()
}
