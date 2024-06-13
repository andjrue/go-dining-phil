package main

import (
	"fmt"
	"sync"
	"time"
)

const numPhil int = 5
const food = 500
const maxSleepTime = 1

var philosophers = [5]int{1, 2, 3, 4, 5}

type Philosopher struct {
	id                  int
	leftFork, rightFork chan struct{}
}

func (p *Philosopher) think() {
	fmt.Printf("Philosopher %d is thinking", p.id)
	time.Sleep(time.Duration(maxSleepTime) * time.Second)
}

func getForks(leftFork, rightFork chan struct{}) bool {
	select {
	case leftFork <- struct{}{}:
		select {
		case rightFork <- struct{}{}:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func putForks(leftFork, rightFork chan struct{}) {
	<-leftFork
	<-rightFork
}

func (p *Philosopher) eatFood(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if getForks(p.leftFork, p.rightFork) {
			fmt.Printf("Philosopher %d is eating.\n", p.id)
			time.Sleep(time.Duration(maxSleepTime) * time.Second)
			p.think()
			return
		}
	}
}

func main() {
	for i := 0; i < numPhil; i++ {
		fmt.Printf("Philosopher #%d\n", philosophers[i])
	}
}
