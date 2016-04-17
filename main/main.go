package main

import (
	ph "dp/philosophers"
	f  "dp/forks"
)

func main() {
	var p [f.FORK_CNT]ph.Philosopher
	for i := 0; i < f.FORK_CNT; i++ {
		p[i] = ph.Philosopher{i, i, (i+1)%f.FORK_CNT}
		go p[i].Run()
	}
	for {}
}
