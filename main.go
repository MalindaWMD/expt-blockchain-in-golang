package main

import (
	"fmt"
	"log"

	"github.com/MalindaWMD/expt-blockchain-in-golang/internal"
)

type App struct {
	bc   *internal.Blockchain
	pool *internal.Mempool
}

func main() {
	bc := internal.NewBlockchain()
	defer bc.DB.Close()

	mempool := internal.NewMempool()

	app := &App{
		bc:   bc,
		pool: mempool,
	}

	from := "1kyGzAzYRdtfvmzo7xnr7wkjEVfDgEwudmovJ8TUaQpi14saATEfueW3dQUkLyGw8owPfsG6"
	to := "121nTVYnA5gaSrTmU1upgkfjEJqYECZYKCfAANgA5U7ZoZGx3F8RhcWAmTjmkfj2CiJMpbbAC"

	// // 1st time only
	// ctx := bc.NewCoinbaseTransaction(from)
	// bc.AddBlock([]*internal.Transaction{ctx})

	fmt.Printf("\n%s balance: %d\n", from, app.balance(from))
	fmt.Printf("%s balance: %d\n\n", to, app.balance(to))

	// FROM 1 to 2
	app.send(from, to, 2)

	// // FROM 2 to 1
	// app.send(to, from, 2)

	fmt.Printf("\n%s balance: %d\n", from, app.balance(from))
	fmt.Printf("%s balance: %d\n\n", to, app.balance(to))

	bc.Print()
}

func (app *App) send(from, to string, amount int) *internal.Transaction {
	tx, err := app.bc.NewTransaction(from, to, amount)
	if err != nil {
		log.Println("TX:", err)
	}

	// Add the mempool
	app.pool.Add(tx)

	fmt.Println("Mempool size: ", len(app.pool.Transactions))

	// TODO: Let miners handle this
	if tx != nil {
		app.bc.AddBlock([]*internal.Transaction{tx})
	}

	return tx
}

func (app *App) balance(address string) int {
	return app.bc.GetBalance(address)
}
