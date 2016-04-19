package philosopher

import (
	"fmt"
	"time"
	"sync"
	"math/rand"
	f "dp/forks"
)

const(
	EAT_MIN_NS   = 1e9  // minimum eat time, ns
	EAT_MAX_NS   = 2e9  // maximum eat time, ns
	THINK_MIN_NS = 5e9  // minimum eat time, ns
	THINK_MAX_NS = 10e9 // maximum eat time, ns
)

type Philosopher struct {
	Id int
	Left_idx int
	Right_idx int
}

var rgen *rand.Rand
var randMutex sync.Mutex

// init initializes the random number generator
func init() {
	rgen = rand.New(rand.NewSource(time.Now().Unix() ) )
	randMutex   = sync.Mutex{}
}

// Think causes a philosopher to sleep for THINK_MIN_NS to THINK_MAX_NS
// nanoseconds, and returns once he's done.
func (p Philosopher) Think() {
	randMutex.Lock()
	r := rgen.Int63n(THINK_MIN_NS) + THINK_MAX_NS - THINK_MIN_NS
	randMutex.Unlock()

	fmt.Printf("[+] Philosopher %d is thinking\r\n", p.Id)
	time.Sleep(time.Duration(r)*time.Nanosecond)
}

// Eat causes a philosopher to eat for EAT_MIN_NS to EAT_MAX_NS nanoseconds,
// and returns once he's done
func (p Philosopher) Eat() {
	randMutex.Lock()
	r := rgen.Int63n(EAT_MIN_NS) + EAT_MAX_NS - EAT_MIN_NS
	randMutex.Unlock()

	fmt.Printf("[+] Philosopher %d is eating\r\n", p.Id)
	time.Sleep(time.Duration(r)*time.Nanosecond)
}

// GetForks acquires forks (i.e., data in the forks channel) from left to
// right. The method uses a monitor to guarantee that philosophers won't be
// preempted between acquiring getting forks, preventing potential deadlock. 
func (p Philosopher) GetForks() {
	f.ForkMutex.L.Lock()
	for {
		select {
		case _ = <-f.Forks[p.Left_idx]:
			_ = <-f.Forks[p.Right_idx]
			f.ForkMutex.L.Unlock()
			return
		default:
			f.ForkMutex.Wait()
		}
	}
}

// ReplaceForks puts single data values back in appropriate forks channels,
// freeing the channels for other philosophers to access and acquire the
// forks
func (p Philosopher) ReplaceForks() {
	f.Forks[p.Left_idx]  <- true
	f.Forks[p.Right_idx] <- true
	f.ForkMutex.Signal()
}

// Run follows the step specifications of the Dining Philosophers problem,
// i.e. Think -> Pick up forks -> Eat -> Put forks back -> Think -> ...
func(p Philosopher) Run() {
	for { 
		p.Think()
		p.GetForks()
		p.Eat()
		p.ReplaceForks()
	}
}
