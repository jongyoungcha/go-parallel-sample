package main

import (
	"fmt"
	"sync"
)

func producer(wg *sync.WaitGroup, queue chan int, trigger chan interface{}) {
	defer wg.Done()
	<-trigger
	for i := 1; i <= 10000; i++ {
		queue <- i
	}
}

func consumer(wg *sync.WaitGroup, queue chan int, trigger chan interface{}) {
	defer wg.Done()
	<-trigger
	for i := range queue {
		fmt.Println("i=", i)
		if i == 10000 {
			break
		}
	}

	fmt.Println("Consumer done!")
}

func main() {
	queue := make(chan int, 10)
	trigger := make(chan interface{})

	wg := &sync.WaitGroup{}

	// make producer
	//for i:=1; i<10; i++ {
	wg.Add(1)
	go producer(wg, queue, trigger)
	//}

	// make consumer
	wg.Add(1)
	go consumer(wg, queue, trigger)

	// start
	//trigger <- struct{}{}
	close(trigger)

	// wait
	wg.Wait()

	close(queue)
	//close(trigger)

	fmt.Println("Done...")
}
