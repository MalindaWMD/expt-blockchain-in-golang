package internal

type Mempool struct {
	Transactions map[string]*Transaction
}

func NewMempool() *Mempool {
	tx := make(map[string]*Transaction)
	return &Mempool{
		Transactions: tx,
	}
}

func (m *Mempool) Add(tx *Transaction) {
	m.Transactions[tx.StringId()] = tx
}

func (m *Mempool) Remove(id string) {
	delete(m.Transactions, id)
}
