package miner

type Miner struct{}

func (m Miner) Mine(txns []string) {
	// try to solve nonce for given txns
	// use channels? run continuously until solved or given new input via a channel
	// return only when solved
	// called as goroutine?
}
