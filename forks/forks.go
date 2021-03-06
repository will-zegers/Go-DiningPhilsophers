package forks

import "sync"

const(
	FORK_CNT = 4 // number of forks/philosophers
)

var ForkMutex sync.Cond
var Forks [FORK_CNT]chan bool

// init initializes the channels that will represent the forks (aggregated as
// an array)
func init() {
	ForkMutex = sync.Cond{L: &sync.Mutex{} }
	for i := 0; i < FORK_CNT; i++ {
		Forks[i] = make(chan bool, 1)
		Forks[i] <- true
	}
}
