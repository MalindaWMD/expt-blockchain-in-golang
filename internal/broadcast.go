package internal

type Broadcaster struct {
	channel chan []string
}

func initBroadcaster() *Broadcaster {
	return &Broadcaster{
		channel: make(chan []string),
	}
}

func (b *Broadcaster) Broadcast(data []string) {
	b.channel <- data
}

func (b *Broadcaster) Listen() <-chan []string {
	return b.channel
}
