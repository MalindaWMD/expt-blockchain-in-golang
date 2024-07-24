package main

import (
	"fmt"
	"log"

	"github.com/MalindaWMD/expt-blockchain-in-golang/internal"
)

func main() {
	bc := internal.NewBlockchain()
	defer bc.DB.Close()

	from := "1PBa8u8mHGHnz1ksCRLX7NsmZTj2uq5htsWTiqoWKU3SRJPiKzQeiH4KL6HmLVp5JDsjymDU"
	to := "1MZurNJ3htDMv4WUxZmbg3cUrkBXRn67XkoXdaC1M7Qp7DcpJBGX6icR58mvt9xZKsgEgpJ8"

	// // 1st time only
	// ctx := bc.NewCoinbaseTransaction(from)
	// bc.AddBlock([]*internal.Transaction{ctx})

	fmt.Printf("\n%s balance: %d\n\n", from, bc.GetBalance(from))
	fmt.Printf("%s balance: %d\n\n", to, bc.GetBalance(to))

	// FROM 1 to 2
	// tx, err := bc.NewTransaction(from, to, 2)
	// if err != nil {
	// 	log.Println("TX:", err)
	// }
	// if tx != nil {
	// 	bc.AddBlock([]*internal.Transaction{tx})
	// }

	// // FROM 2 to 1
	tx, err := bc.NewTransaction(to, from, 2)
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

	fmt.Printf("\n%s balance: %d\n", from, bc.GetBalance(from))
	fmt.Printf("%s balance: %d\n\n", to, bc.GetBalance(to))

	bc.Print()
}
