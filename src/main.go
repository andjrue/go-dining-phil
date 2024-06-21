package main

import (
	"fmt"
	"sync"
	"time"
)

// Constants

const numPhil int = 200
const maxSleepTime int = 10

// Variables to be updated later

var food int = 500000
var foodMutex sync.Mutex
var m = make(map[int]int)

type Philosopher struct {
	id        int
	leftFork  *sync.Mutex
	rightFork *sync.Mutex
}

func (p *Philosopher) think() {
	fmt.Printf("Philosopher %d is thinking \n", p.id)
	time.Sleep(time.Duration(maxSleepTime) * time.Millisecond)
}

func (p *Philosopher) eatFood(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		foodMutex.Lock()
		p.rightFork.Lock()
		p.leftFork.Lock()

		if food > 0 {
			fmt.Printf("Philosopher %d is eating \n", p.id)
			food--
			m[p.id]++ // Increase food eaten
			fmt.Printf("Food Remaining: %d \n", food)
			foodMutex.Unlock()

			time.Sleep(time.Duration(maxSleepTime) * time.Millisecond)
			p.think()

		} else {
			p.rightFork.Unlock()
			p.leftFork.Unlock()
			foodMutex.Unlock()
			return
		}

		// After we eat we can release the forks
		p.rightFork.Unlock()
		p.leftFork.Unlock()
	}
}

func createPhilArr(numPhil int) []*Philosopher {
	forks := make([]*sync.Mutex, numPhil)
	for i := 0; i < numPhil; i++ {
		forks[i] = &sync.Mutex{}
	}

	philosophers := make([]*Philosopher, numPhil)
	for i := 0; i < numPhil; i++ {
		philosophers[i] = &Philosopher{
			id:        i + 1,
			leftFork:  forks[i],
			rightFork: forks[(i+1)%numPhil],
		}
	}

	return philosophers
}

func main() {
	var wg sync.WaitGroup

	philosophers := createPhilArr(numPhil)

	for _, philosopher := range philosophers {
		wg.Add(1)
		go philosopher.eatFood(&wg)
	}

	wg.Wait()
	fmt.Println("All food has been eaten.")
	fmt.Printf("Food Map: \n %d", m)
}
