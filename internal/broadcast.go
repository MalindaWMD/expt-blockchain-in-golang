package internal

type Broadcaster struct {
	txChannel    chan []string
	blockChannel chan string
}

func initBroadcaster() *Broadcaster {
	return &Broadcaster{
		txChannel:    make(chan []string),
		blockChannel: make(chan string),
	}
}

func (b *Broadcaster) BroadcastTransaction(data []string) {
	b.txChannel <- data
}

func (b *Broadcaster) ListenTransaction() <-chan []string {
	return b.txChannel
}

func (b *Broadcaster) BroadcastBlock(data string) {
	b.blockChannel <- data
}

func (b *Broadcaster) ListenBlock() <-chan string {
	return b.blockChannel
}
