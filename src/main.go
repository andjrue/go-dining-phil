package main

import (
	"fmt"
	"sync"
	"time"
)

const numPhil int = 10
const maxSleepTime = 50

var food = 500
var foodMutex sync.Mutex
var m = make(map[int]int)

type Philosopher struct {
	id                  int
	leftFork, rightFork *sync.Mutex
	mu                  sync.Mutex
}

func (p *Philosopher) think() {
	fmt.Printf("Philosopher %d is thinking \n", p.id)
	time.Sleep(time.Duration(maxSleepTime) * time.Millisecond)
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
		foodMutex.Lock()
		p.leftFork.Lock()
		p.rightFork.Lock()

		if food > 0 {
			fmt.Printf("Philosopher %d is eating. \n", p.id)
			food--
			m[p.id]++
			fmt.Printf("Food Remaining: %d \n", food)
			foodMutex.Unlock()

			time.Sleep(time.Duration(maxSleepTime) * time.Millisecond)
			p.think()
		} else {
			foodMutex.Unlock()
			p.rightFork.Unlock()
			p.leftFork.Unlock()
			return
		}

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
	fmt.Println("All food has been eaten. Map is below.\n")
	fmt.Printf("Food Map: \n %d", m)

}
