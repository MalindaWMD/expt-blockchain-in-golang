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

func (m *Mempool) Get(count int) []*Transaction {
	var txs []*Transaction
	itemCount := 0
	for _, value := range m.Transactions {
		if itemCount > count {
			break
		}
		txs = append(txs, value)
	}

	return txs
}
