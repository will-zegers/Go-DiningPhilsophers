package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)

const(
	PHIL_CNT     = 4    // number of philosophers
	EAT_MIN_NS   = 1e9  // minimum eat time, ns
	EAT_MAX_NS   = 2e9  // maximum eat time, ns
	THINK_MIN_NS = 5e9  // minimum eat time, ns
	THINK_MAX_NS = 10e9 // maximum eat time, ns
)

type Philosopher struct {
	id int
	left_idx int
	right_idx int
}

var rgen *rand.Rand
var forks [PHIL_CNT]chan bool
var wg sync.WaitGroup

// init initializes the random number generator, as well as the channels that
// will represent the forks (aggregated as an array)
func init() {
	rgen = rand.New(rand.NewSource(time.Now().Unix() ) )
	for i := 0; i < PHIL_CNT; i++ {
		forks[i] = make(chan bool, 1)
		forks[i] <- true
	}
}

// Think causes a philosopher to sleep for THINK_MIN_NS to THINK_MAX_NS
// nanoseconds, and returns once he's done.
func (p Philosopher) Think() {
	r := rgen.Int63n(THINK_MIN_NS) + THINK_MAX_NS - THINK_MIN_NS
	fmt.Printf("Philosopher %d is thinking\r\n", p.id)
	time.Sleep(time.Duration(r)*time.Nanosecond)
}

// Eat causes a philosopher to eat for EAT_MIN_NS to EAT_MAX_NS nanoseconds,
// and returns once he's done
func (p Philosopher) Eat() {
	r := rgen.Int63n(EAT_MIN_NS) + EAT_MAX_NS - EAT_MIN_NS
	fmt.Printf("Philosopher %d is eating\r\n", p.id)
	time.Sleep(time.Duration(r)*time.Nanosecond)
}

// GetForks acquires forks (i.e., data in the forks channel) from left to
// right. After a call to get forks, the accessed channels will be empty,
// causing other philosophers attempting to access these channels to block
func (p Philosopher) GetForks() {
	_ = <-forks[p.left_idx]
	_ = <-forks[p.right_idx]
}

// ReplaceForks puts single data values back in appropriate forks channels,
// freeing the channels for other philosophers to access and acquire the
// forks
func (p Philosopher) ReplaceForks() {
	forks[p.left_idx]  <- true
	forks[p.right_idx] <- true
}

// Run follows the step specifications of the Dining Philosophers problem,
// i.e. Think -> Pick up forks -> Eat -> Put forks back -> Think -> ...
func(p Philosopher) Run() {
	for i := 0; i < 5; i++ {
		p.Think()
		p.GetForks()
		p.Eat()
		p.ReplaceForks()
	}
	fmt.Printf("Philosopher %d has finished\r\n", p.id)
	wg.Done()
}

func main() {
	var p [PHIL_CNT]Philosopher
	for i := 0; i < PHIL_CNT; i++ {
		p[i] = Philosopher{i, i, (i+1)%PHIL_CNT}
		go p[i].Run()
		wg.Add(1)
	}
	wg.Wait()
}
