package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func producer(start, step, end int, inCh chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start; i <= end; i += step {
		time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
		inCh <- i
	}
}

func consumer(inCh <-chan int, outCh chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range inCh {
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
		outCh <- val * val
	}
}

func finalFilter(outCh <-chan int, done chan<- struct{}) {
	defer close(done)
	var lastPrinted *int
	for val := range outCh {
		if lastPrinted == nil || val > *lastPrinted {
			fmt.Println(val)
			lastPrinted = &val
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	inCh := make(chan int, 5)
	outCh := make(chan int, 5)

	var prodWg sync.WaitGroup
	var consWg sync.WaitGroup
	done := make(chan struct{})

	// Stage 1: Producers
	prodWg.Add(2)
	go producer(1, 2, 29, inCh, &prodWg)
	go producer(2, 2, 30, inCh, &prodWg)

	go func() {
		prodWg.Wait()
		close(inCh)
	}()

	// Stage 2: Consumers
	consWg.Add(2)
	go consumer(inCh, outCh, &consWg)
	go consumer(inCh, outCh, &consWg)

	go func() {
		consWg.Wait()
		close(outCh)
	}()

	// Stage 3: Final Filter
	go finalFilter(outCh, done)

	<-done
}
