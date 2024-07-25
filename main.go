package main

import (
	"fmt"

	"github.com/MalindaWMD/expt-blockchain-in-golang/internal"
)

func main() {
	bc := internal.NewBlockchain()
	defer bc.DB.Close()

	from := "1kyGzAzYRdtfvmzo7xnr7wkjEVfDgEwudmovJ8TUaQpi14saATEfueW3dQUkLyGw8owPfsG6"
	to := "121nTVYnA5gaSrTmU1upgkfjEJqYECZYKCfAANgA5U7ZoZGx3F8RhcWAmTjmkfj2CiJMpbbAC"

	// // 1st time only
	// ctx := bc.NewCoinbaseTransaction(from)
	// bc.AddBlock([]*internal.Transaction{ctx})

	fmt.Printf("\n%s balance: %d\n", from, bc.GetBalance(from))
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
	// tx, err := bc.NewTransaction(to, from, 2)
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
